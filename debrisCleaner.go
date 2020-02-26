/* Copyright 2020 by n3o33 <discord n3o33#2384>
 * Not proprietary and confidential, 
 * feel free to share, copy and change
 */

 /*
 * VERSION 1.0
 */

  /* DESCRIPTION
 This script will keep own debris fields clean
 */

func collectDebris() {
	slots = GetSlots()
	freeSlots = slots.Total - slots.InUse
    if freeSlots > GetFleetSlotsReserved() {
		celestials,_ = GetCelestials()
		for celestial in celestials {
			systemInfos, _ = GalaxyInfos(celestial.GetCoordinate().Galaxy, celestial.GetCoordinate().System)
			planetInfo = systemInfos.Position(celestial.GetCoordinate().Position)
			if planetInfo != nil {
				needRecs = planetInfo.Debris.RecyclersNeeded
				for needRecs > 0 {
					ships, _ = celestial.GetShips()
					fleet = NewFleet()
					fleet.SetOrigin(celestial.GetCoordinate())
					fleet.SetDestination(NewCoordinate(celestial.GetCoordinate().Galaxy, celestial.GetCoordinate().System,celestial.GetCoordinate().Position, DEBRIS_TYPE))
					fleet.SetMission(RECYCLEDEBRISFIELD)
					fleet.SetSpeed(HUNDRED_PERCENT)
					if ships.Recycler < needRecs {
						fleet.AddShips(RECYCLER, ships.Recycler)
					} else {
						fleet.AddShips(RECYCLER, needRecs)
					}
					f, err = fleet.SendNow()
					LogInfo("[DEBRIS] send fleet " + f)
					if err != nil {
						LogError("[DEBRIS] send fleet error with " + err)
						break
					} else if ships.Recycler < needRecs {
						needRecs -= ships.Recycler
						Sleep((f.BackIn + 5) * 1000)
					} else {
						break
					}
				}
			}
		}
	}
}

collectDebris()
