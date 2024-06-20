#include <sys/epoll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>
#include <stdio.h>

#define MAXEVENTS 128
#define SERV_PORT 8080

int make_nonblocking(int fd) {
    int flags = fcntl(fd, F_GETFL, 0);
    if (flags == -1) {
        return -1;
    }
    flags |= O_NONBLOCK;
    if (fcntl(fd, F_SETFL, flags) == -1) {
        return -1;
    }
    return 0;
}

//创建一个listen的fd，设置为非阻塞的
int tcp_nonblocking_server_listen(int port) {
    int listen_fd;
    struct sockaddr_in server_addr;
    listen_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd == -1) {
        perror("socket");
        return -1;
    }
    make_nonblocking(listen_fd);

    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(port);

    if (bind(listen_fd, (struct sockaddr *) &server_addr, sizeof(server_addr)) == -1) {
        perror("bind");
        close(listen_fd);
        return -1;
    }

    if (listen(listen_fd, SOMAXCONN) == -1) {
        perror("listen");
        close(listen_fd);
        return -1;
    }

    return listen_fd;
}

int main() {
    int listen_fd, socket_fd;
    int n, i;
    int efd;
    struct epoll_event event;
    struct epoll_event *events;
    char buf[512];


    listen_fd = tcp_nonblocking_server_listen(SERV_PORT); // 创建一个listen
    if (listen_fd == -1) {
        return -1;
    }

    efd = epoll_create1(0); // 创建epoll
    if (efd == -1) {
        perror("epoll_create1");
        close(listen_fd);
        return -1;
    }

    event.data.fd = listen_fd; // 设置监听事件
    event.events = EPOLLIN | EPOLLET; // 设置可读写
    if (epoll_ctl(efd, EPOLL_CTL_ADD, listen_fd, &event) == -1) { // 加入到epoll中
        perror("epoll_ctl");
        close(listen_fd);
        close(efd);
        return -1;
    }

    events = calloc(MAXEVENTS, sizeof(event)); // 给数组分配内存
    if (!events) {
        perror("calloc");
        close(listen_fd);
        close(efd);
        return -1;
    }

    while (1) {
        n = epoll_wait(efd, events, MAXEVENTS, -1); // 等待感兴趣

        for (i = 0; i < n; i++) {
            if ((events[i].events & EPOLLERR) || // 事件出错
                (events[i].events & EPOLLHUP) ||
                !(events[i].events & EPOLLIN)) {
                close(events[i].data.fd);
                continue;
            }

            if (listen_fd == events[i].data.fd) { // 事件id == 监听id
                struct sockaddr_storage ss;
                socklen_t slen = sizeof(ss);
                int fd = accept(listen_fd, (struct sockaddr *) &ss, &slen);
                if (fd == -1) {
                    if (errno != EAGAIN && errno != EWOULDBLOCK) {
                        perror("accept");
                    }
                    continue;
                }

                make_nonblocking(fd); // 创建非阻塞io连接建立
                event.data.fd = fd;
                event.events = EPOLLIN | EPOLLET;
                if (epoll_ctl(efd, EPOLL_CTL_ADD, fd, &event) == -1) { // 把连接加入到epoll
                    perror("epoll_ctl");
                    close(fd);
                }
                continue;
            }

            socket_fd = events[i].data.fd;
            n = read(socket_fd, buf, sizeof(buf)); // 读取read数据
            if (n <= 0) {
                if (n == 0 || (n < 0 && errno != EAGAIN)) {
                    close(socket_fd); // 出错关闭或连接关闭
                }
            } else {
                write(socket_fd, buf, n); // 写请求
                // 处理其他逻辑
            }
        }
    }

    free(events);
    close(listen_fd);
    close(efd);
    return 0;
}
