#include <poll.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <errno.h>
#include <stdio.h>
#include <string.h>

#define SERV_PORT 8080
#define MAX_CLIENTS 1024
#define BUF_SIZE 512

// 设置文件描述符为非阻塞模式
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

// 创建并绑定一个非阻塞的服务器监听套接字
int tcp_nonblocking_server_listen(int port) {
    int listen_fd;
    struct sockaddr_in server_addr;

    // 创建套接字
    listen_fd = socket(AF_INET, SOCK_STREAM, 0);
    if (listen_fd == -1) {
        perror("socket");
        return -1;
    }

    // 设置套接字为非阻塞模式
    if (make_nonblocking(listen_fd) == -1) {
        perror("make_nonblocking");
        close(listen_fd);
        return -1;
    }

    // 初始化服务器地址结构
    server_addr.sin_family = AF_INET;
    server_addr.sin_addr.s_addr = INADDR_ANY;
    server_addr.sin_port = htons(port);

    // 绑定套接字到指定端口
    if (bind(listen_fd, (struct sockaddr *) &server_addr, sizeof(server_addr)) == -1) {
        perror("bind");
        close(listen_fd);
        return -1;
    }

    // 开始监听
    if (listen(listen_fd, SOMAXCONN) == -1) {
        perror("listen");
        close(listen_fd);
        return -1;
    }

    return listen_fd;
}

int main() {
    int listen_fd, conn_fd, sock_fd;
    struct pollfd fds[MAX_CLIENTS];
    int nfds = 1;
    struct sockaddr_in client_addr;
    socklen_t client_len;
    char buf[BUF_SIZE];
    int n, i;

    // 创建并绑定监听套接字
    listen_fd = tcp_nonblocking_server_listen(SERV_PORT);
    if (listen_fd == -1) {
        return -1;
    }

    // 初始化 poll 文件描述符结构
    memset(fds, 0, sizeof(fds));
    fds[0].fd = listen_fd;
    fds[0].events = POLLIN;

    while (1) {
        // 调用 poll 函数监控文件描述符集合中的事件
        int ready = poll(fds, nfds, -1);
        if (ready == -1) {
            perror("poll");
            break;
        }

        // 检查监听套接字是否有新的连接请求
        if (fds[0].revents & POLLIN) {
            client_len = sizeof(client_addr);
            conn_fd = accept(listen_fd, (struct sockaddr *) &client_addr, &client_len);
            if (conn_fd == -1) {
                perror("accept");
                continue;
            }

            // 设置新连接的套接字为非阻塞模式
            if (make_nonblocking(conn_fd) == -1) {
                perror("make_nonblocking");
                close(conn_fd);
                continue;
            }

            // 将新连接的套接字加入 poll 文件描述符集合
            fds[nfds].fd = conn_fd;
            fds[nfds].events = POLLIN;
            nfds++;
        }

        // 遍历所有可能的文件描述符，检查是否有可读事件
        for (i = 1; i < nfds; i++) {
            if (fds[i].revents & POLLIN) {
                sock_fd = fds[i].fd;
                n = read(sock_fd, buf, sizeof(buf));
                if (n <= 0) {
                    if (n == 0 || (n < 0 && errno != EAGAIN)) {
                        // 连接关闭或读取出错，关闭套接字并从集合中移除
                        close(sock_fd);
                        fds[i] = fds[nfds - 1];
                        nfds--;
                        i--;
                    }
                } else {
                    // 将读取到的数据回写给客户端
                    write(sock_fd, buf, n);
                    // 处理其他逻辑
                }
            }
        }
    }

    // 关闭监听套接字
    close(listen_fd);
    return 0;
}