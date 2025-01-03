## ⚙️ Configuration

The Elvia CLI support a configuration file (`3lv.yml`) to store some default values for the commands.
Using a configuration file means you can avoid typing the same options every time you run a command.
The configuration file is optional and should be placed in the root of your repository.

Example (see [pkg/shared/config.go](https://github.com/3lvia/cli/blob/trunk/pkg/shared/config.go) for the full specification):

```yaml
# 3lv.yml
system: core
applications:
  - name: demo-api
    projectFile: applications/DemoApi/DemoApi.csproj
    helmValuesFile: .github/deploy/values-demo-api.yml
    buildContext: . # optional
```

With the above configuration, you can for example run the `build` command using just the application name:

```bash
3lv build demo-api
```
