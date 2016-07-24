package dotnet

type Runtime struct {}

func NewRuntime() (err error, runtime Runtime) {
  runtime = Runtime{}
  return err, runtime
}
