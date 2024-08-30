### 寄存器
容量非常小，通常是32位、64位或更少的几个寄存器。寄存器专门用于执行CPU指令. 中央处理器（CPU）内部的一部分，存储指令执行过程中需要的操作数、地址或结果，以及控制信息。CPU首先从寄存器中获取数据，若寄存器中没有，CPU会从L1缓存中查找，如果L1缓存也没有，继续查找L2缓存，依此类推，最后才访问主存
* 通用寄存器：用于存储一般数据，比如操作数或计算结果。
* 指令寄存器（IR）：用于存储当前正在执行的指令。
* 程序计数器（PC）：用于存储下一条将要执行的指令的地址。
* 状态寄存器：存储运算结果的状态信息，如零标志、进位标志等。




### 内存
通过内存地址(指针)来访问真实内存的数据
物理地址空间 : 能索引到内存单元的地址合集(其实就是一个这样的整数所表示的范围)
虚拟地址空间 : 它不能真正索引到具体的设备单元,需要一个MMU（内存管理单元）虚拟地址转换成真正的物理地址才能访问相应的设备
* 避免修改到其他程序的内存数据，所以提供了有一个命名空间的概念。每个程序只能操作自己的命名空间下的内存（都是起始从0开始）

#### 内存区域
* 代码区
    * 在运行时会把代码编译的指令存放在这里（只读的无法修改）

* 数据区 (可读写)
    * 已初始化全局变量及静态变量  Data
    * 未初始化的全局变量和静态变量 Bss

* 常量区 (只读: const)
    * 字符串字面量和其他常量

* 栈 （高效内存管理）—— 数据暂时存储的地方
    * 程序执行过程中自动管理的内存区域，主要存储局部变量和函数调用信息。内存分配是有编译器或处理器自动管理
    * 栈的分配和释放操作非常简单（只需调整栈指针）
    * 栈中的数据缓存到 L1 或 L2 缓存中提高访问速度

* 堆 （灵活内存管理——不受栈空间限制）
    * 程序执行过程中动态分配的，用于分配对象或数据结构的
    * 需要维护一个自由内存块的列表，以便找到合适的块来满足内存分配请求.操作涉及查找和合并空闲块，因此访问速度相对较慢。
    * 访问模式通常是随机的，且不会像栈那样有很好的局部性。


#### 堆栈总结
堆栈(堆段和栈段的大小都是动态增加和减少的、且增长方向相反。堆是向高地址方向增长，栈是向低地址方向增长)
由于cup缓存的空间的问题，不可能把所有数据都放入到cpu缓存中，那么如果把经常用的放到cpu缓存中，不常用的放到内存缓存中，这样就可以加快速度，限制如下
1. 由于CPU缓存空间有限，所以通常存放经常访问的小数据（如基本数据类型或指针），特别是具有良好局部性的数据（如局部变量、函数参数、返回地址等上下文信息）。CPU缓存的使用主要依赖于数据的访问模式，由硬件自动管理。
2. 为了快速分配和释放内存，函数调用时在现有的栈上“分配”一个新的栈帧，这只是调整栈指针的位置。栈帧的销毁也是通过栈指针的移动自动实现的，这种自动化管理简化了编程工作。
3. 由于作用域的限制，函数返回后栈帧会被销毁，因此如果某些数据需要在函数返回后继续使用，通常会将其复制到堆中。堆不仅用于存放较大的内存对象，更常用于动态分配生命周期较长的数据，特别是需要在多个函数之间共享的数据或复杂的数据结构。
4. 堆中的内容如果没有引用，通常不会立即被删除，而是由垃圾回收机制来处理