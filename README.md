# go-blog

学了一段时间的GO语言了，期间断断续续，但对GO语言的敬爱从未减退。因为比较熟悉WEB编程,所以就想写个博客之类的东西玩一玩。
由于学得不好，又没有使用框架，所以写得比较渣。

## 热编译

写这种服务器和应用一体的程序如果没有热编译是非常痛苦的，这里使用了beego框架的bee工具。

https://github.com/beego/bee

    $ go get https://github.com/beego/bee
    $ cd your/proj/path
    $ bee run
    
## 数据库

使用了sqlite3数据库，挺好用的。

http://www.sqlite.org

位置：  ./db/blog.sq3

## 编辑器

使用了百度的UEditor

http://ueditor.baidu.com/website

## 后台

路径：  localhost:8080/admin

账号：  admin

密码：  1992
