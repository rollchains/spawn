

spawn module new nameserviceibc --ibc-middleware


app.go:

    (so is spawn an --ibc-middleware or ibc.Module?)

	var nameserviceStack porttypes.IBCModule
	nameserviceStack = nameserviceibc.NewIBCMiddleware(transferStack, app.NameserviceibcKeeper)


	ibcRouter.AddRoute(nameserviceibctypes.ModuleName, nameserviceStack)
	app.IBCKeeper.SetRouter(ibcRouter)


proto/nameserviceibc/genesis.proto


	  string port_id = 2;
