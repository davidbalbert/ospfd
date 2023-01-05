package ospf

type Protocol struct {
	RouterID RouterID
	Areas    map[AreaID]Area
	// TODO: VirtualLinks
	// TODO: ExternalRoutes
	// TODO: LSDB (or maybe just AS external?)
	// TODO: RIB

	Config Config
}
