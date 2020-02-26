STRING = import("strings")
buildList = []
func checkCoordinate(c, field) {
    coord = STRING.Replace(c, "[", "", -1)
    coord = STRING.Replace(coord, "]", "", -1)
    coord, err = ParseCoord(coord)
    if err != nil {
        LogError("unrecognized err", field, err)
        StopScript(__FILE__) }
        return coord
}
func checkBuildingName(s, field) {
    if !( s == "Metal" || s == "Crystal" || s == "Deuterium") {
        LogError("unrecognized building name", field)
        StopScript(__FILE__)
    }
}
func checkIfNumber(s, field) {
    if Atoi(s) <= 0 {
        LogError("unrecognized number", field);
        StopScript(__FILE__)
    }
}
func checkResearchName(s, field) {
    if !( s == "PlasmaTechnology" || s == "Astrophysics") {
        LogError("unrecognized research name", field)
        StopScript(__FILE__)
    }
}
func parseBuildInfo(bi) {
    a = []
    for line in STRING.Split(bi, "\n") {
        if len(line) > 0 {
            field = STRING.Fields(line)
            if Atoi(field[0]) > 0 {
                if len(field) == 5 {
                    coord = checkCoordinate(field[2], field)
                    checkBuildingName(field[3], field)
                    checkIfNumber(field[4], field)
                    a += { coord  : { field[3] : field[4] } }
                } else if len(field) == 4 {
                    checkResearchName(field[1]+field[2], field)
                    checkIfNumber(field[3], field)
                    a += { "0:0:0" : { field[1]+field[2] : field[3] } }
                } else if len(field) == 3 {
                    checkResearchName(field[1], field)
                    if STRING.Contains(field[2], "+") {
                        for sub in STRING.Split(field[2], "+") {
                            a += { "0:0:0" : { field[1] : sub } }
                        }
                    } else {
                        checkIfNumber(field[2], field)
                        a += { "0:0:0" : { field[1] : field[2] } }
                    }
                } else {
                    LogError("unrecognized", field)
                    StopScript(__FILE__)
                }
            }
        }
    }
    return a
}
func buildFromList(bl) {
    for entry in bl {
        for coord, tobuild in entry {
            for what, lvl in tobuild {
                infoComptBuild(coord, what, lvl)
            }
        }
    }
}
func transformInfoComptToNinja(s) {
    switch s {
        case "PlasmaTechnology": return PLASMATECHNOLOGY
        case "Astrophysics": return ASTROPHYSICS
        case "Metal": return METALMINE
        case "Crystal": return CRYSTALMINE
        case "Deuterium": return DEUTERIUMSYNTHESIZER
        default: LogError("unrecognized", s)
        StopScript(__FILE__)
    }
}
func infoComptBuild(c, w, l) {
    if c == "0:0:0" {
        LogInfo("add to queue", c, w, l)
        //err = AddItemToQueue(GetCachedCelestial(NewCoordinate(GetHomeWorld().GetCoordinate().Galaxy, GetHomeWorld().GetCoordinate().System, GetHomeWorld().GetCoordinate().Position, PLANET_TYPE)).ID, transformInfoComptToNinja(w), 0)
        //if err != nil {
        //    LogError("err adding to queu", err)
        //}
    } else {
        LogInfo("add to queue", c, w, l)
        //err = AddItemToQueue(GetCachedCelestial(c).ID, transformInfoComptToNinja(w), 0)
        //if err != nil {
        //    LogError("err adding to queu", err)
        //}
    }
}

// start
buildFromList(parseBuildInfo(infoCompt))
