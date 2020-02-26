/* Copyright 2020 by n3o33 <discord n3o33#2384>
 * Not proprietary and confidential, 
 * feel free to share, copy and change
 */

 /*
 * VERSION 1.0
 */

  /* DESCRIPTION
 This script will repatriate resources from planet to homeworld if a storage is full
 */

func repatriate(res, from) {
    myShips, _ = from.GetShips()
    LogInfo("[STORAGE] resources to send " + res)
    pf, lc, sc, cargo = CalcFastCargoPF(myShips.Pathfinder, myShips.LargeCargo, myShips.SmallCargo, res.Total())
    fleet = NewFleet()
    fleet.SetOrigin(from.GetCoordinate())
    fleet.SetDestination(GetHomeWorld().GetCoordinate())
    fleet.SetMission(TRANSPORT)
    fleet.SetSpeed(HUNDRED_PERCENT)
	fleet.SetResources(res)
	fleet.AddShips(PATHFINDER, pf)
    fleet.AddShips(LARGECARGO, lc)
    fleet.AddShips(SMALLCARGO, sc)
    f, err = fleet.SendNow()
    LogInfo("[STORAGE] send fleet " + f)
    if err != nil {
        LogError("[STORAGE] send fleet error with " + err)
    }
}

func checkIfStorageIsFull() {
 	slots = GetSlots()
	freeSlots = slots.Total - slots.InUse
    if freeSlots > GetFleetSlotsReserved() {
		celestials,_ = GetCelestials()
		for celestial in celestials {
			if celestial.GetType() == PLANET_TYPE {
				resDetails, _ = celestial.GetResourcesDetails()
				resToSend = NewResources(0,0,0)
				LogInfo("[STORAGE] " + celestial.GetCoordinate() + " " + celestial.GetName())
				if resDetails.Metal.Available >= resDetails.Metal.StorageCapacity { 
					LogInfo("[STORAGE] ---- metal full")
					if sendAllRes {
						resToSend = NewResources(resDetails.Metal.Available, resDetails.Crystal.Available, resDetails.Deuterium.Available - leaveDeuterium)
					} else {
						resToSend = resToSend.Add(NewResources(resDetails.Metal.Available, 0, 0))
					}
				} else {
					LogInfo("[STORAGE] ---- metal storage is " + Round(((100/resDetails.Metal.StorageCapacity)*resDetails.Metal.Available)) + " %")
				}
				if resDetails.Crystal.Available >= resDetails.Crystal.StorageCapacity { 
					LogInfo("[STORAGE] ---- crystal full")
					if sendAllRes {
						resToSend = NewResources(resDetails.Metal.Available, resDetails.Crystal.Available, resDetails.Deuterium.Available - leaveDeuterium)
					} else {
						resToSend = resToSend.Add(NewResources(0, resDetails.Crystal.Available, 0))
					}
				} else {
					LogInfo("[STORAGE] ---- crystal storage is " + Round(((100/resDetails.Crystal.StorageCapacity)*resDetails.Crystal.Available)) + " %")
				}
				if resDetails.Deuterium.Available >= resDetails.Deuterium.StorageCapacity { 
					LogInfo("[STORAGE] ---- deuterium full")
					if sendAllRes {
						resToSend = NewResources(resDetails.Metal.Available, resDetails.Crystal.Available, resDetails.Deuterium.Available - leaveDeuterium)
					} else {
						resToSend = resToSend.Add(NewResources(0, 0, resDetails.Deuterium.Available))
					}
				} else {
					LogInfo("[STORAGE] ---- deuterium storage is " + Round(((100/resDetails.Deuterium.StorageCapacity)*resDetails.Deuterium.Available)) + " %")
				}
				if resToSend.Total() > 0 {
				   repatriate(resToSend, celestial) 
				}
			}
		}
	}
}