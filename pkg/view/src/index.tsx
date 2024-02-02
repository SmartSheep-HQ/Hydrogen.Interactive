import "solid-devtools";

/* @refresh reload */
import { render } from "solid-js/web";

import "./index.css";
import "./assets/fonts/fonts.css";
import { lazy } from "solid-js";
import { Route, Router } from "@solidjs/router";

import "@fortawesome/fontawesome-free/css/all.css";

import RootLayout from "./layouts/RootLayout.tsx";
import { UserinfoProvider } from "./stores/userinfo.tsx";
import { WellKnownProvider } from "./stores/wellKnown.tsx";

const root = document.getElementById("root");

render(() => (
  <WellKnownProvider>
    <UserinfoProvider>
      <Router root={RootLayout}>
        <Route path="/" component={lazy(() => import("./pages/feed.tsx"))} />
        <Route path="/auth" component={lazy(() => import("./pages/auth/callout.tsx"))} />
        <Route path="/auth/callback" component={lazy(() => import("./pages/auth/callback.tsx"))} />
      </Router>
    </UserinfoProvider>
  </WellKnownProvider>
), root!);
