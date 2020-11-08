package dto

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

//定义一个全局存在的哈希环，并用init()方法初始化
var GlobalHashCircle *Consistent

func init() {
	GlobalHashCircle = NewConsistent()
}

//声明新切片类型，并定义其用于排序的三个函数
type sortedHash []uint32

func (sh sortedHash) Len() int {
	return len(sh)
}

func (sh sortedHash) Less(i, j int) bool {
	return sh[i] < sh[j]
}

func (sh sortedHash) Swap(i, j int) {
	sh[i], sh[j] = sh[j], sh[i]
}

//创建结构体保存一致性hash信息
type Consistent struct {
	Circle       map[uint32]string //hash环，key为哈希值，值存放结点的信息
	SortedHashes sortedHash        //已经排好序的circle中Hash值
	VirtualNode  int               //虚拟结点个数，用来增加hash的平衡性
	sync.RWMutex                   //map 读写锁
}

//上面那个结构体的构造函数，设置默认结点数量
func NewConsistent() *Consistent {
	return &Consistent{
		//初始化变量
		Circle: make(map[uint32]string),
		//设置虚拟结点个数
		VirtualNode: 20,
	}
}

//向hash环中添加结点
func (c *Consistent) Add(serverInfo string) {
	c.Lock()
	defer c.Unlock()
	//根据生成的结点以及虚拟结点添加到hash环中
	for i := 0; i < c.VirtualNode; i++ {
		c.Circle[c.GetHashKey(serverInfo, i)] = serverInfo
	}
	//更新排序
	c.updateSortedHashes()
}

//根据一个IP地址+端口，从hash环中删除对应的结点
func (c *Consistent) Remove(serverInfo string) {
	c.Lock()
	defer c.Unlock()
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.Circle, c.GetHashKey(serverInfo, i))
	}
	c.updateSortedHashes()
}

//根据一个string，获得对应的hash值，并获取hash环中第一个在其前面的hash值对应的IP地址+端口（二分查找）
func (c *Consistent) Get(data string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.Circle) == 0 {
		return "", errors.New("Hash circle is empty! ")
	}
	hash := c.GetHashKey(data, 0)
	serverInfo := c.search(hash)
	return serverInfo, nil
}

//根据接收的info和虚拟结点的index，构造其在hash环上的位置（也就是其hash值）
func (c *Consistent) GetHashKey(info string, index int) uint32 {
	key := info + strconv.Itoa(index)
	//使用IEEE多项式返回数据的CRC-32校验和
	return crc32.ChecksumIEEE([]byte(key))
}

//对Hash环中数据按照Hash大小依次排序，并将排序结果存在sortedHashes中
func (c *Consistent) updateSortedHashes() {
	var hash sortedHash
	for k := range c.Circle {
		hash = append(hash, k)
	}
	sort.Sort(hash)
	c.SortedHashes = hash
}

//用二分法搜寻第一个在此hash值前面的hash值对应的IP+端口
func (c *Consistent) search(hash uint32) string {
	if hash <= c.SortedHashes[0] || hash > c.SortedHashes[len(c.SortedHashes)-1] {
		return c.Circle[c.SortedHashes[0]]
	}
	l := 0
	r := len(c.SortedHashes)-1
	res := 0
	for l < r {
		mid := (l + r) / 2
		if c.SortedHashes[mid] >= hash && c.SortedHashes[mid-1] < hash {
			res = mid
			break
		} else if c.SortedHashes[mid] < hash {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}
	return c.Circle[c.SortedHashes[res]]
}