package landingparser

import (

)

type LandingParser struct {
}

// #\sTitle\n+([\s]?.*)\n+---\n+#\sDescription\n+([\s]*.+)+---\n+#\sImages\n(\*\s!\[\]\((.+/.+)\))