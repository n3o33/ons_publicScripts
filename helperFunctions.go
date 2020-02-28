/* Copyright 2020 by n3o33 <discord n3o33#2384>
 * Not proprietary and confidential, 
 * feel free to share, copy and change
 */

 /*
 * VERSION 1.0
 */

/* DESCRIPTION
 * Helpfull functions for all ! 
 * Support me to get a whole llist of formula implementations ^^ 
*/

func getMetalProduction(lvlMetalMine, lvlPlasmaTechnology) { 
    return Floor(30 * (1 + (lvlPlasmaTechnology / 100)) * GetUniverseSpeed() * lvlMetalMine * Pow(1.1, lvlMetalMine) + (30 * GetUniverseSpeed())) 
}

func getCrystalProduction(lvlCrystalMine, lvlPlasmaTechnology) {
    return Floor((15 * GetUniverseSpeed()) + (20 * GetUniverseSpeed() * (1 + (lvlPlasmaTechnology * 0.0066)) * lvlCrystalMine * Pow(1.1, lvlCrystalMine)))
}

func getDeuteriumProduction(celestial, lvlDeuteriumSynthesizer, lvlFusionReactor, lvlPlasmaTechnology) {
    deutProd = GetUniverseSpeed() * (10 * lvlDeuteriumSynthesizer * Pow(1.1, lvlDeuteriumSynthesizer) * (1.36 - (0.004 * celestial.Temperature.Mean())))
    if lvlFusionReactor > 0 {
        fusionCons = Floor(-10 * GetUniverseSpeed() * lvlFusionReactor * Pow(1.1, lvlFusionReactor))
        deutProd = Floor(deutProd + fusionCons) * (1 + (lvlPlasmaTechnology * 0.0033))
    } else {
        deutProd = Floor(deutProd * (1 + (lvlPlasmaTechnology * 0.0033)))
    }    
}

func getEnergyProduction(celestial, lvlFusionReactor, lvlEnergyTechnology, lvlSolarPlant, solarSatellite) {
    eFusion = 30 * lvlFusionReactor * Pow((1.05 + lvlEnergyTechnology * 0.01), lvlFusionReactor)
    ePlant = 20 * lvlSolarPlant * Pow(1.1, lvlSolarPlant)
    eSolSat = ((celestial.Temperature.Mean() + 160) / 6) * solarSatellite
    return Floor(eFusion + ePlant + eSolSat)
}

func getEnergyAvailble(celestial, lvlFusionReactor, lvlEnergyTechnology, lvlSolarPlant, solarSatellite, lvlMetalMine, lvlCrystalMine, lvlDeuteriumSynthesizer) {
    metalMine = 10 * lvlMetalMine * Pow(1.1, lvlMetalMine)
    crystalMine = 10 * lvlCrystalMine * Pow(1.1, lvlCrystalMine)
    deutSynth = 20 * lvlDeuteriumSynthesizer * Pow(1.1, lvlDeuteriumSynthesizer)
    return getEnergyProduction(celestial, lvlFusionReactor, lvlEnergyTechnology, lvlSolarPlant, solarSatellite) - Floor(metalMine + crystalMine + deutSynth)
}

func getMetalStorageCapaticy(lvlMetalStorage) {
    E = 2.718281828459045
    return (5000 * Floor((2.5 * Pow(E, 20 / 33 * lvlMetalStorage))))
    }
}

func getCrystalStorageCapaticy(lvlCrystalStorage) {
    E = 2.718281828459045
    return (5000 * Floor((2.5 * Pow(E, 20 / 33 * lvlCrystalStorage))))
    }
}

func getMetalStorageCapaticy(lvlDeuteriumStorage) {
    E = 2.718281828459045
    return (5000 * Floor((2.5 * Pow(E, 20 / 33 * lvlDeuteriumStorage))))
    }
}

func deleteElementFromArray(array, elementToDelete) {
    tmp = []
    for e in array { 
        if e != elementToDelete {
            tmp += e 
        } 
    }
    return tmp 
}
a = [1,2,3,4,5]
Print(a, "after", deleteElementFromArray(a, 3))

func getBuildTime(buildingID, lvlBuilding, resourceCost, lvlRoboterFactoryplanet, lvlShipyard, lvlNaniteFactory) {
    if buildingID == NANITEFACTORY {
        return ((resourceCost.Metal + resourceCost.Crystal) / (2500 * (1 + lvlRoboterFactoryplanet) * Pow(2, lvlNaniteFactory) * GetUniverseSpeed())) * 60 * 60
    } else {
        return ((resourceCost.Metal + resourceCost.Crystal) / (2500 * Max(4 - lvlBuilding / 2, 1) * (1 + lvlRoboterFactoryplanet) * Pow(2, lvlNaniteFactory) * GetUniverseSpeed())) * 60 * 60
    }
}
Print("TEST getBuildTime for NANI"      , getBuildTime(NANITEFACTORY, 6, NewResources(32000000, 16000000, 3200000), 10, 12, 5))
Print("TEST getBuildTime for MetalMine" , getBuildTime(MetalMine, 26, NewResources(1515000, 378767, 0), 10, 12, 5))