build:
    go build -o slack-tableflip .

package:
    ko build --local .
