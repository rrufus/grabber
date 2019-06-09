# Grab

grabber.Grab is a function to carry out multiple http GETs in one grab (i.e. parallel).

It will return a map containing all of the results once the last GET has completed.

It will return an error if one of the http GETs returns an error.

All you need to supply is an http client (http.DefaultClient will do), and an Item struct per call, where items includes the `Key` the result is stored under in the resulting map, and the `Url` to call.