rust中方法能作为入参吗
结论
Rust 里方法可以作为函数入参，但要区分：关联方法（静态方法） 和 实例方法，二者写法不一样。
1. 关联方法（static 方法）
不带 &self，本质就是普通函数，可以直接当作参数传递。
rust
运行
struct Foo;

impl Foo {
    // 关联方法
    fn func(x: i32) -> i32 {
        x * 2
    }
}

// 接收一个函数类型作为参数
fn call(f: fn(i32) -> i32, val: i32) -> i32 {
    f(val)
}

fn main() {
    // 直接把关联方法当作参数传入
    let res = call(Foo::func, 10);
    println!("{}", res); // 20
}
函数类型：fn(i32) -> i32。
2. 实例方法（带 &self / &mut self）
实例方法类型会自带隐式第一个参数 self，不能直接传给普通函数指针。
写法 1：取方法指针（未绑定实例）
rust
运行
struct Foo {
    num: i32,
}

impl Foo {
    fn method(&self, x: i32) -> i32 {
        self.num + x
    }
}

fn main() {
    // 未绑定实例的方法指针：fn(&Foo, i32) -> i32
    let m: fn(&Foo, i32) -> i32 = Foo::method;

    let foo = Foo { num: 5 };
    println!("{}", m(&foo, 3)); // 8
}
类型展开：
fn(&Self, 其他参数) -> 返回值。
写法 2：闭包绑定实例（最常用）
把实例捕获进闭包，再把闭包传进去：
rust
运行
fn run<F: Fn(i32) -> i32>(f: F, arg: i32) -> i32 {
    f(arg)
}

fn main() {
    let foo = Foo { num: 5 };
    // 闭包捕获 foo
    let closure = |x| foo.method(x);
    let res = run(closure, 4);
    println!("{}", res);
}
3. 三种可接收的参数形式
函数指针 fn(...) -> T：只能传普通函数 + 关联方法 + 未绑定实例的方法指针，无捕获；
Fn(...) -> T trait：可以传函数指针、闭包（可捕获环境），优先用这个；
FnMut / FnOnce：对应可变捕获、只执行一次的闭包。
4. 关键总结
静态关联方法：可以直接作为 fn 类型参数；
实例方法：
不取实例：写成 Type::method，类型为 fn(&Type, args...)；
要固定实例：包一层闭包再传入；
业务代码优先用泛型约束 F: Fn(...) -> T，兼容性最强。