In Go, a build tag, or a build constraint, is an identifier added to a piece of code that determines when the file should be included in a package during the build process. This allows to build different versions of Go application from the same source code and to toggle between them in a fast and organized manner. Many developers use build tags to improve the workflow of building cross-platform compatible applications, such as programs that require code changes to account for variances between different operating systems. Build tags are also used for integration testing, allowing you to quickly switch between the integrated code and the code with a mock service or stub, and for differing levels of feature sets within an application.

Let’s take the problem of differing customer feature sets as an example. When writing some applications, you may want to control which features to include in the binary, such as an application that offers Free, Pro, and Enterprise levels. As the customer increases their subscription level in these applications, more features become unlocked and available. To solve this problem, you could maintain separate projects and try to keep them in sync with each other through the use of import statements. While this approach would work, over time it would become tedious and error prone. An alternative approach would be to use build tags.

In this article, we will use build tags in Go to generate different executable binaries that offer Free, Pro, and Enterprise feature sets of a sample application. Each will have a different set of features available, with the Free version being the default.

## Building the Free Version

Let’s start by building the Free version of the application, as it will be the default when running **`go build` without any build tags**. Later on, we will use build tags to selectively add other parts to the program.

```go
package main

import "fmt"

var features = []string{
  "Free Feature #1",
  "Free Feature #2",
}

func main() {
  for _, f := range features {
    fmt.Println(">", f)
  }
}
```
In the `main.go` file, we created a program that declares a slice named `features`, which holds two strings that represent the features of our Free application. The *`main()`* function uses a `for loop` to range through the `features slice` and print all of the features available to the screen.

After build and run the program we’ll receive the following output:
```
> Free Feature #1
> Free Feature #2
```
The program has printed out two free features, completing the Free version of the app.

So far, we created an application that has a very basic feature set. Next, we will build a way to add more features into the application at build time without modifying the main file.

## Adding the Pro Features With `go build`
Since we can’t edit the `main.go` file, we’ll need to use another mechanism for injecting more `features` into the `features slice` using ***build tags***.

Let’s create a new file called `pro.go` that will use an *`init()`* function to append more features to the features slice:

```go
package main

func init() {
  features = append(features,
    "Pro Feature #1",
    "Pro Feature #2",
  )
}
```

In this code, we used *`init()`* to run code before the *`main()`* function of our application, followed by *`append()`* to add the Pro features to the features slice.

After build and run the program we’ll receive the following output:

```
> Free Feature #1
> Free Feature #2
> Pro Feature #1
> Pro Feature #2
```
The application now includes both the Pro and the Free features. However, this is not desirable: since there is no distinction between versions, the Free version now includes the features that are supposed to be only available in the Pro version. To fix this, you could include more code to manage the different tiers of the application, or you could use build tags to tell the Go tool chain which .go files to build and which to ignore. Let’s add build tags in the next step.

## Adding Build Tags
Now we can use build tags to distinguish the Pro version of the application from the Free version.

A build tag looks like:
``` go
//go:build tag_name
```

By putting this line of code as the first line of the package and replacing **`tag_name`** with the name of desired build tag, you will tag this package as code that can be selectively included in the final binary. Let’s see this in action by adding a build tag to the `pro.go` file to tell the **go build** command to ignore it unless the tag is specified.

```go
//go:build pro

package main

func init() {
  features = append(features,
    "Pro Feature #1",
    "Pro Feature #2",
  )
}
```

At the top of the `pro.go` file, the *`//go:build pro`* was added followed by a blank newline. This trailing newline is required, otherwise Go interprets this as a comment. Build tag declarations must also be at the very top of a `.go` file. 
> **Nothing, not even comments, can be above build tags**.

The *`go:build pro`* declaration tells the `go build` command that this isn’t a comment, but instead is a **build tag**. The second part is the pro tag. By adding this tag at the top of the pro.go file, the go build command will now only include the pro.go file with the pro tag is present.

Compile and run the application again, it will produce following output:
```
> Free Feature #1
> Free Feature #2
```

Since the `pro.go` file requires a *pro tag* to be present, the file is ignored and the application compiles without it.

When running the `go build` command, we can use the **-tags** flag to conditionally include code in the compiled source by adding the tag itself as an argument. Let’s do this for the pro tag:

```
go build -tags pro
```
Compilled app will output the following:

```
> Free Feature #1
> Free Feature #2
> Pro Feature #1
> Pro Feature #2
```
Now we only get the extra features when we build the application using the pro build tag.

This is fine if there are only two versions, but things get complicated when you add in more tags. To add in the Enterprise version of the app in the next step, we will use multiple build tags joined together with Boolean logic.

## Build Tag Boolean Logic
When there are multiple build tags in a Go package, the tags interact with each other using Boolean logic. To demonstrate this, we will add the Enterprise level of our application using both the *pro tag* and the *enterprise tag*.

In order to build an Enterprise binary, we will need to include both the default features, the Pro level features, and a new set of features for Enterprise. First, we need to create new file that will add the new Enterprise features, called `enterprise.go` with the following content:
```go
package main

func init() {
  features = append(features,
    "Enterprise Feature #1",
    "Enterprise Feature #2",
  )
}
```
Currently this file doesn't have any build tags, this means that these features will be added to the Free version when executing `go build`. For `pro.go`, we added `//go:build pro` to the top of the file to tell `go build` that it should only be included when `-tags pro` is used. In this situation, you only needed one build tag to accomplish the goal. When adding the new Enterprise features, however, you first must also have the Pro features.

If we add *`//go:build enterprise`* to the top of the file, as we do in the `pro.go`, and build it, the Pro features will be lost. The solution is the following:
```go
//go:build pro && enterprise
```

And build stage will be
```
go build -tags pro,enterprise
```
We need to provide both tags to include Pro and Enterprise features.

More about boolean logic in build constraints can be found in [official documentation](https://pkg.go.dev/cmd/go#hdr-Build_constraints).

# Links
- Original [article](https://www.digitalocean.com/community/tutorials/customizing-go-binaries-with-build-tags) from Digital Ocean
