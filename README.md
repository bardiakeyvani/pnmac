# What does this program do?
1. Read weathers table and returns the day number with  smallest temperature spread
2. Read soccers table and returns the name of the team with smallest For and Against score

## required argument
- type: --weather or --soccer
- table: path to table file. Example: --table=./tables/soccer.dat

# How to run
You can run the included binary on any macOS operating system. You can also build it from the source code on any OS with [golang](https://golang.org/doc/install) installed.
## Using binary 

```
./bin/pnmac --soccer --table=./tables/soccer.dat

./bin/pnmac --weather --table=./tables/w_data.dat 
```

## build from source
```
git clone git@github.com:bardiakeyvani/pnmac.git
cd pnmac
go build -o pnmac 
pnmac --weather --table=./tables/w_data.dat
pnmac --soccer --table=./tables/soccer.dat
```
