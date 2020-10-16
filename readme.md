#### servdown

basic uptime monitoring tool

```sh
git clone git@github.com:demon36/servdown.git
cd servdown
go build
./servdown #this will emit a sample json file
# nano servdown.json # edit config
# cp ./servdown ~/bin/ # copy to "PATH" included path
# echo "servdown" | at now # run as daemon
```