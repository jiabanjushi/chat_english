package tools

import (
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"
	"testing"
)

func TestMyTest(t *testing.T) {
	MyTest()
}
func TestMyStruct(t *testing.T) {
	MyStruct()
}
func TestMySlice(t *testing.T) {
	//数组
	arr := [4]int{1, 2, 3, 4}
	arrSlice := arr[:]
	log.Printf("长度:%d,容量:%d\n", len(arrSlice), cap(arrSlice))
	//容量超过 ,扩容问2倍 4*2=8
	arrNewSlice := append(arrSlice, 5)
	log.Printf("长度:%d,容量:%d\n", len(arrNewSlice), cap(arrNewSlice))

	var arr2 = [1024]int{}
	for i := 0; i < 1024; i++ {
		arr2[i] = i
	}
	log.Printf("长度:%d,容量:%d\n", len(arr2[:]), cap(arr2[:]))
	//容量超过1024 ,扩容问1.25倍 1024*1.25=1280
	arr2NewSlice := append(arr2[:], 1024)
	log.Printf("长度:%d,容量:%d\n", len(arr2NewSlice), cap(arr2NewSlice))

	arr3NewSlice := make([]int, 2, 1024)
	arr3NewSlice[0] = 1
	log.Printf("%v,长度:%d,容量:%d\n", arr3NewSlice, len(arr3NewSlice), cap(arr3NewSlice))
}
func TestMyInt(t *testing.T) {
	var a int8
	a = 127
	log.Println(a)
}

type Test struct {
	Name string
}

func f(out io.Writer) {
	// ...do something...
	if out != nil {
		out.Write([]byte("done!\n"))
	}
}
func TestMyInterface(t *testing.T) {
	//var buf *bytes.Buffer
	//f(buf)
	//arr := []int{1, 2}
	//testInterface2(arr)
	//a := (*interface{})(nil)
	//log.Printf("%v\n", a == nil)
	//var c interface{}
	//c = (*interface{})(nil)
	//log.Printf("%v\n", c == nil)
	//
	//b := &Test{
	//	Name: "sss",
	//}
	//testInterface(b)
}
func testInterface2(b interface{}) {
	t := reflect.TypeOf(b).String()
	fmt.Printf("%#v", t)
	c := b.([]int)
	for _, d := range c {
		fmt.Printf("%#v", d)
	}
}

func testInterface(b interface{}) {
	b.(*Test).Name = "aaa"
	fmt.Println(b.(*Test).Name)
}

type SMap struct {
	sync.RWMutex
	Map map[int]int
}

func (l *SMap) readMap(key int) (int, bool) {
	l.RLock()
	value, ok := l.Map[key]
	l.RUnlock()
	return value, ok
}

func (l *SMap) writeMap(key int, value int) {
	l.Lock()
	l.Map[key] = value
	l.Unlock()
}

var mMap *SMap

func TestMyMap(t *testing.T) {
	mMap = &SMap{
		Map: make(map[int]int),
	}

	for i := 0; i < 10000; i++ {
		go func() {
			mMap.writeMap(i, i)
		}()
		go readMap(i)
	}
}
func readMap(i int) (int, bool) {
	return mMap.readMap(i)
}

type Ages []int

func (a *Ages) AgeAdd(num int) {
	ages := append([]int(*a), num)
	*a = ages
}

func TestMyage(t *testing.T) {
	var ages Ages
	ages.AgeAdd(10)
	ages.AgeAdd(20)
	t.Logf("%+v", ages)
}
