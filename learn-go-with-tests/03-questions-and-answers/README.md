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

## [Error types](https://quii.gitbook.io/learn-go-with-tests/questions-and-answers/error-types)

***Creating your own types for errors can be an elegant way of tidying up your code, making your code easier to use and test.***


Before:

- code
    ```go
    if res.StatusCode != http.StatusOK {
        return "", fmt.Errorf("did not get 200 from %s, got %d", url, res.StatusCode)
    }
    ```
- test
    ```go
	want := fmt.Sprintf("did not get 200 from %s, got %d", svr.URL, http.StatusTeapot)
	got := err.Error()

	if got != want {
		t.Errorf(`got "%v", want "%v"`, got, want)
	}
    ```
- problems
    - same string in prod code and test codes
    - annoying to read and write
    - exact error message is not what we're concerned with

After:
- Use custom Error type

    ```go
    type BadStatusError struct {
        URL    string
        Status int
    }

    func (b BadStatusError) Error() string {
        return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
    }
    ```
- code
    ```go
	if res.StatusCode != http.StatusOK {
		return "", BadStatusError{URL: url, Status: res.StatusCode}
	}
    ```
- test

    ```go
    var got BadStatusError
    isBadStatusError := errors.As(err, &got)
    want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

	if !isBadStatusError {
		t.Fatalf("was not a BadStatusError, got %T", err)
	}
    ```
- improvements
    - `DumbGetter` gets simpler
    - Enable more sophisticated error handling with a type assertion
    - Still an `error`. we can treat it in the same way as other errors.
