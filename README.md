To build:
```
make
```

To install:
```
sudo make install
```

`make` will build a binary called metro for you which will allow you to search routes, stops and departures. To start off with you can run `./metro routes` which will list all the routes provided by svc.metrotransit.org. You can then pass an argument to the command to filter your search request down. For instances `./metro routes target` should list a handfull of items. If you pass the --direction flag you will see the directions this route goes. `./metro routes target --direction`

Using the stops subcommand you can see the locations the route stops on. You must pass a route and a direction substring. These substrings must match only one item or you will see a list of all items matching your substring. `./metro stops "harmar target" south` will list all stops on the HarMar Target - Lexington Av Route going south.

Lastly the departures subcommand will list all the departures for a stop on a given route with a specified direction. `./metro departures blue south TF22` will list all departures on the Blue Line, heading south from the Target Field Station Platform 2 stop. If you specify the `--next` flag it will only show you the next departure. Lastly there is a --template flag that allows you to pass in a custom golang template. To only see the next departure time you would use the following command `metro departures blue south TF22 --next --template '{{ .DepartureText }}'`
