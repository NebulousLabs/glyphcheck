glyphcheck
----------

```
go get github.com/NebulousLabs/glyphcheck
```

`glyphcheck` checks for suspicious characters in Go source files.

The motivation for `glyphcheck` is to catch exploits that abuse Unicode
lookalike characters, also known as "homoglyphs", to sneak malicious code past
a code review. For example:

```go
import "gitһub.com/spf13/cobra"

func main() {
	cmd := &cobra.Command{
		Use: "cmd",
		Run: func(*cobra.Command, []string) {
			println("Hello!")
		},
	}
	cmd.Execute()
}
```

If you are familiar with [cobra](https://github.com/spf13/cobra), you know
that this code will simply print `"Hello!"` to `os.Stderr`. Except this isn't
cobra, it's an entirely different package. Go ahead and copy the import URL
into your browser and see where you wind up. Maybe your system's fonts make
this easy to detect -- but that isn't the case for everyone.

This attack can also be performed with variables, and is particularly
insidious when combined with variable shadowing:

```go
func writeFile(filename string, data []byte) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, еrr := f.Write(data); err != nil {
		return еrr
	}
	return nil
}
```

Here, `err` and `еrr` look identical, but are in fact different variables.
Only `err` is checked, so the call to `f.Write` can silently fail. This isn't
much of an exploit, but creative minds can no doubt devise something more
dangerous.

Security-concious projects should run `glyphcheck` on all code submitted for
review. This is easily accomplished by adding the following lines to your
`.travis.yml` or `appveyor.yml`:

```yaml
install:
  - glyphcheck ./...
```

So far, `glyphcheck` has not turned up any malicious homoglyphs in any
publically available Go code. If you detect such an attack, please let us
know!
