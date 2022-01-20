# Questions and answers

## OS Exec

Condition:
- `os/exec.Command()` is called in `GetData()`.
- Want to use test data in tests.
Question: Should I add a "test" flag like `GetData(mode string)`?
Answer: [Dependency Injection](https://quii.gitbook.io/learn-go-with-tests/go-fundamentals/dependency-injection)

1. before

    ```go
	cmd := exec.Command("cat", "msg.xml")

	out, _ := cmd.StdoutPipe()
	var payload Payload
	decoder := xml.NewDecoder(out)
    ```

1. after

    ```go
    func GetData(data io.Reader) string {
        var payload Payload
        xml.NewDecoder(data).Decode(&payload)
        return strings.ToUpper(payload.Message)
    }

    func getXMLFromCommand() io.Reader {
        cmd := exec.Command("cat", "msg.xml")
        out, _ := cmd.StdoutPipe()

        cmd.Start()
        data, _ := ioutil.ReadAll(out)
        cmd.Wait()

        return bytes.NewReader(data)
    }
    ```
1. Now we can test `GetData`.

    ```go
    func TestGetData(t *testing.T) {
        input := strings.NewReader(`
    <payload>
        <message>Cats are the best animal</message>
    </payload>`)

        got := GetData(input)
        want := "CATS ARE THE BEST ANIMAL"

        if got != want {
            t.Errorf("got %q, want %q", got, want)
        }
    }
    ```
