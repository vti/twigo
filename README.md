A static blog engine written in Go for educational purposes. It is mostly
a rewrite of http://github.com/vti/Twist, a blog engine written in Perl.

    git clone https://github.com/vti/twigo

Get `goop` if you don't have it:

    go get github.com/nitrous-io/goop

Install dependencies

    goop install

Start twigo

    goop exec ./twigo serve --conf conf.json --listen :8080
