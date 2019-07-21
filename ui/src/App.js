import React from "react";
import "./App.css";
import SelectTicker from "./SelectTicker";
import {BrowserView, isBrowser, isMobile, MobileView} from "react-device-detect";

function App() {
    if (isBrowser) {
        return (
            <BrowserView>
                <div className="App">
                    <header className="container-columns-header">Market Patterns</header>
                    <SelectTicker/>
                </div>
            </BrowserView>
        );
    } else if (isMobile) {
        return (
            <MobileView>
                <div className="App">
                    <header className="container-columns-header">Market Patterns</header>
                    <SelectTicker/>
                </div>
            </MobileView>
        );
    }

    return (<div>Unknown view display, not a browser or a device</div>);
}

export default App;