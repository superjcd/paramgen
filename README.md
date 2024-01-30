# paramgen 
`paramgen` is a  tiny code gen tool eyes on reducing the time of writing boilerplate codes.  
given a struct like:  
```golang
type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
	IsAdmin  int
}
```

we want something like:  
```golang
type UserForm struct {
	Name     string `form:"name"`
	Password string `form:"password"`
	Email    string `form:"email"`
	IsAdmin  int    `form:"is_admin"`
}

type UserJson struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	IsAdmin  int    `json:"is_admin"`
}

```
this is particularly useful when working with [gin](https://github.com/gin-gonic/gin) or other kinds of web framework demanding some sort of tag annotation.

## How to use
To intall `paramgen`:
```shell
go install github.com/superjcd/paramgen
```

Put a notation at the top of model  declaration file, just like:
```
//go:generate  paramgen

package model

type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
	IsAdmin  int
}
```
Final step is just cast the spell `go generate ./...` , you shall see the auto-generated file lies beside the initial struct declaration files:

## How it works
For those curious, you can read the blog linked just below  

## Reference
[the-ultimate-guide-to-writing-a-go-tool](https://arslan.io/2017/09/14/the-ultimate-guide-to-writing-a-go-tool/)