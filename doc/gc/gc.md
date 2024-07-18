### 垃圾回收机制
* 引用计数器：每个对象都为何一个引用计数，被引用时+1，引用解除时-1，直到0时销毁 （可以很快回收无需等到垃圾回收触发时但是需要实时维护） ———— swift
  * 垃圾回收可以迅速回收，无需等待垃圾回收触发 （单一处理）
* 标记清除：从根遍历所有引用对象，如果没有被标记的视为清除 （无需实时维护但是需要STW时间导致整个不可用）———— go
  * 垃圾回收等待某个节点触发，需要暂停业务扫描清除。（批量处理）
* 分代收集：按照对象的生命周期时长划分不同的代空间，不同代有不同的回收算法和回收频率(性能好可以对不同的代进行配置) ———— java
  * 可以使用不同的回收算法，提高效率


### STW
1. 暂停所有业务逻辑 
2. 开始标记找出能可达/不可达对象 
3. 标记完成清除不可达对象

### 三色标记法 
* 黑色，已经扫描过的对象+子对象为灰色； 
* 灰色，已经扫描过的对象； 
* 白色，没有扫描的对象

1. 所有对象刚开始都是白色，标识没有遍历过的数据
2. GC开始时扫描root set标记为灰色，如果发现子对象指针全部扫描到（灰色）就改为黑色
3. 扫描全部后剩下的白色节点就是没有已经没被引用，放到freelist，则可以被清除

满足强弱之一，即可保证对象不丢失
* 强三色不变式 ： 强制不允许黑色对象引用白色对象 
* 弱三色不变式：黑色可以引用白色对象，但需要白色对象存在其他的灰色对象对他引用

###  go 1.1版本 （该算法是最简单的，但是整个流程业务不可用，可能到秒级别的）
STW stop the world 停止所有业务
Mark 通过root set找出可达对象进行标记
Sweep 把未标记的对象加入freelist
Start the world 开始所有业务

### go 1.3 Mark与Sweep分离，Sweep可不在STW中
多协程同时mark减少了STW的时间
把Sweep后台与程序运行同时并行（已经不到达的放到freelist的，期间不可能对他们再有操作）

### go 1.5 三色标记法
1. 所有的对象都是白色对象
2. 多协程从roots set一直遍历扫描,当扫描到标色为灰色对象，一直遍历扫描，如果当前对象的子对象全部都扫描过（灰色），此时该对象标识为黑色。
3. 如何判断扫描是否结束，只要有灰色标记就代表还有子对象没被处理完。（引用计数器的思路，每处理完一个--1直到0，即可sweep）
(因为不会STW的话，在对数据的删除或插入的时候，需要引入写屏障/删除屏障，通知垃圾回收变化———— 把自己的父对象改为灰色，自己也是灰色)

###  屏障
* 插入写屏障: A对象引用B对象，B对象被标记灰色，如果没有写屏障此时B依旧是白色的就可能被清除 
* 删除写屏障: 被删除的对象，如果自身为灰色或者白色，那么被标记为灰色，等到下次GC再处理（弱三色不变式 保护灰色对象到白色对象的路径不会断）


优化方案：
本来整个过程本来是需要暂停程序不可用的，为了减少程序不可用的时间，优化如下
1. 在标记过程中是多线程进行处理，
2. 把Mark和Sweep进行分开，Mark是不可用的，sweep是后台慢慢清理的
3. 也把Mark不可用改为可用，由于标记过程中程序可用就有可能数据会修改，那么对修改的数据及对应的父对象修改为黑色及灰色保证本次GC不会受到影响