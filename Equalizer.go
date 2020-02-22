/* Copyright 2020 by n3o33 <discord n3o33#2384>
 * Not proprietary and confidential, 
 * feel free to share, copy and change
 */

//####################### CONFIG #############################
//############################################################
splitPlanetWorkers                  = "allPlanets"
//splitPlanetWorkers                = ["M:1:2:4", "P:2:3:4"]
splitMetal                          = false
splitCrystal                        = false
splitDeuterium                      = true
roundFactor                         = 500000
//############################################################

STRINGS = import("strings")
func sendFleet(from, to, resources, mission, ships) {
    Print("from", from, "to", to, "resources", resources, "mission", mission, "ships", ships)
    fleet = nil
    speed = HUNDRED_PERCENT
    for {
        slots = GetSlots()
        freeSlots = slots.Total - slots.InUse
        if freeSlots <= GetFleetSlotsReserved() {
            fleet, _ = GetFleets()
                waitTime = 0
                for f in fleet {
                    if waitTime == 0 || f.BackIn < waitTime {
                        waitTime = f.BackIn
                    }
                }
                LogInfo("no slots free, wait " + ShortDur(waitTime))
                Sleep((waitTime + 5) * 1000)
        } else {
            f = NewFleet()
            f.SetOrigin(from)
            f.SetDestination(to)
            if resources != nil {
                f.SetResources(resources)
            }
            f.SetSpeed(speed)
            f.SetMission(mission)
            f.SetShips(ships)
            f, err = f.SendNow()
            if err != nil {
                if STRINGS.Contains(err.Error(), "not enough cargo capacity for fuel") {
                    LogError("send fleet error, going to reduce speed")
                    if speed < 1 {
                        LogError("send fleet error on 10% speed")
                        break
                    } else {
                        speed --
                    }
                } else {
                    LogError("send fleet error", err)
                    break
                }
            } else {
                fleet = f
                break
            }
        }
    }
    return fleet
}

func doWork() {
    LogInfo("start ....")
    allRes          = NewResources(0, 0, 0)
    planets         = 0
    celCanOffer     = {}
    celNeed         = {}
    planetList      = []
    for c, r in GetAllResources()[0] {
        allRes = allRes.Add(r)
        celCanOffer[GetCachedCelestial(c).Coordinate] = r
    }
    if splitPlanetWorkers == "allPlanets" {
        for p in GetCachedPlanets() {
            planets ++
            planetList += p
        }
    } else {
        for planet in splitPlanetWorkers {
            planets ++
            p = GetPlanet(planet)[0]; Sleep(Random(500,800))
            planetList += p
        }
    }
    divMetal        = allRes.Metal / planets
    divCrystal      = allRes.Crystal / planets
    divDeuterium    = allRes.Deuterium / planets
    resPerCel       = NewResources(divMetal, divCrystal, divDeuterium)
    for planet in planetList {
        resCel = NewResources(0, 0, 0)
        resCel = resCel.Add(planet.GetResources()[0]); Sleep(Random(500, 800))
        resNeedOrOfferMetal = resCel.Metal - resPerCel.Metal
        resNeedOrOfferCrystal = resCel.Crystal - resPerCel.Crystal
        resNeedOrOfferDeuterium = resCel.Deuterium - resPerCel.Deuterium
        need = NewResources(0, 0, 0)
        if resNeedOrOfferMetal < 0 { need = need.Add(NewResources(Abs(resNeedOrOfferMetal), 0, 0)) }
        if resNeedOrOfferCrystal < 0 { need = need.Add(NewResources(0, Abs(resNeedOrOfferCrystal), 0)) }
        if resNeedOrOfferDeuterium < 0 { need = need.Add(NewResources(0, 0, Abs(resNeedOrOfferDeuterium))) }
        if need.Total() > 0 { celNeed[planet] = need }
    }
    // make final list to send res from
    finalList = []
    for coordNeed, _ in celNeed {
        for start in GetSortedCelestials(coordNeed.Coordinate) {
            if start.GetType() == PLANET_TYPE {
                // check if needer
                for c, _ in celNeed {
                    if c.Coordinate.Equal(start.Coordinate) {
                        celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(resPerCel)
                    }
                }
            }
            if !coordNeed.Coordinate.Equal(start.Coordinate) {
                if celNeed[coordNeed].Total() > 0 {
                    resToSend = NewResources(0, 0, 0)
                    if celNeed[coordNeed].Metal > 0 && splitMetal {
                        if celCanOffer[start.Coordinate].Metal > celNeed[coordNeed].Metal {
                            resToSend = resToSend.Add(NewResources(celNeed[coordNeed].Metal, 0, 0))
                            celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(NewResources(celNeed[coordNeed].Metal, 0, 0))
                            celNeed[coordNeed] = celNeed[coordNeed].Sub(NewResources(celNeed[coordNeed].Metal, 0, 0))
                        } else {
                            resToSend = resToSend.Add(NewResources(celCanOffer[start.Coordinate].Metal, 0, 0))
                            celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(NewResources(celCanOffer[start.Coordinate].Metal, 0, 0))
                            celNeed[coordNeed] = celNeed[coordNeed].Sub(NewResources(celCanOffer[start.Coordinate].Metal, 0, 0))
                        }
                    }
                    if celNeed[coordNeed].Crystal > 0 && splitCrystal {
                        if celCanOffer[start.Coordinate].Crystal > celNeed[coordNeed].Crystal {
                            resToSend = resToSend.Add(NewResources(0, celNeed[coordNeed].Crystal, 0))
                            celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(NewResources(0, celNeed[coordNeed].Crystal, 0))
                            celNeed[coordNeed] = celNeed[coordNeed].Sub(NewResources(0, celNeed[coordNeed].Crystal, 0))
                        } else {
                            resToSend = resToSend.Add(NewResources(0, celCanOffer[start.Coordinate].Crystal, 0))
                            celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(NewResources(0, celCanOffer[start.Coordinate].Crystal, 0))
                            celNeed[coordNeed] = celNeed[coordNeed].Sub(NewResources(0, celCanOffer[start.Coordinate].Crystal, 0))
                        }
                    }
                    if celNeed[coordNeed].Deuterium > 0 && splitDeuterium {
                        if celCanOffer[start.Coordinate].Deuterium > celNeed[coordNeed].Deuterium {
                            resToSend = resToSend.Add(NewResources(0, 0, celNeed[coordNeed].Deuterium))
                            celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(NewResources(0, 0, celNeed[coordNeed].Deuterium))
                            celNeed[coordNeed] = celNeed[coordNeed].Sub(NewResources(0, 0, celNeed[coordNeed].Deuterium))
                        } else {
                            resToSend = resToSend.Add(NewResources(0, 0, celCanOffer[start.Coordinate].Deuterium))
                            celCanOffer[start.Coordinate] = celCanOffer[start.Coordinate].Sub(NewResources(0, 0, celCanOffer[start.Coordinate].Deuterium))
                            celNeed[coordNeed] = celNeed[coordNeed].Sub(NewResources(0, 0, celCanOffer[start.Coordinate].Deuterium))
                        }
                    }
                    resToSend = NewResources(Floor(resToSend.Metal / roundFactor) * roundFactor, Floor(resToSend.Crystal / roundFactor) * roundFactor, Floor(resToSend.Deuterium / roundFactor) * roundFactor)
                    if resToSend.Total() > 0 { finalList += { start : { coordNeed : resToSend } } }
                }
            }
        }
    }
    
    for inner in finalList {
        for from , i in inner {
            for to, res in i {
                LogInfo("send resources", from.Coordinate, "--->", to.Coordinate, "with", res)
                myShips, _ = from.GetShips();  Sleep(Random(500, 800))
                lc, sc, cargo = CalcFastCargo(myShips.LargeCargo, myShips.SmallCargo, res.Total())
                if cargo > 0 {
                    tmpShips = NewShipsInfos()
                    tmpShips.Set(LARGECARGO, lc)
                    tmpShips.Set(SMALLCARGO, sc)
                    f = sendFleet(from.Coordinate, to.Coordinate, res, TRANSPORT, *tmpShips)
                    Sleep(Random(15, 50) * 1000)
                } else {
                    LogError("not enough cargo to send")
                }
            }
        }
    }
    LogInfo(".... done")
}

doWork()