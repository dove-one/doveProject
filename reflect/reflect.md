##### Reflect

```
Kind:Kind是一个枚举值，用来判断操作的对应类型，例如是否是指针，数组等
所以reflect调用的不对，直接就panic

反射输出所有的字段名字，关键点在于，只有Kind==struct的才有字段
注：无指针类型

func (v Value) Field(i int) Value {...}
func (v Value) FieldByName(name string) Value {...}
Field性能优先于FieldByName

用反射设置值
可以用反射来修改一个字段的值。需要注意的是，修改字段的值之前一定要检查CanSet()
简单来说，就是必须使用结构体指针，那么结构体的字段才是可以修改的。
当然指针指向的对象也是可以修改的

TODO：refelct api待补充

反射面试
· 什么是反射？反射可以看做是对对象和对类型的描述，而我们可以通过反射来间接操作对象
· 反射的使用场景？
· 能不能通过反射修改方法？不能。为什么不能？go runtime没有暴露接口
· 什么样的字段可以被反射修改？有一个方法CanSet可以判断，简单来说就是addressable

```

##### 