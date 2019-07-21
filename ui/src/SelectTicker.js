import React from 'react';
import SelectSeries from "./SelectSeries";
import GraphStockPrice from "./GraphStockPrice";
import GraphPatternDensity from "./GraphPatternDensity"
import {BrowserView, isBrowser, MobileView} from "react-device-detect";

class SelectTicker extends React.Component {

    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount() {

        fetch('http://localhost:7666/api/latest/tickers/')
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            this.setState(
                {
                    selectedSymbol: data.names[0].symbol,
                    selectedCompany: data.names[0].company,
                    tickers: data.names
                }
            );
        });
    }

    handleChange(event) {
        let tickers;
        if (this.state != null && this.state.tickers != null) {
            tickers = this.state.tickers;
        }

        this.setState({selectedSymbol: event.target.value, tickers: tickers});
    }

    render() {
        if (this.state == null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.selectedSymbol == null) {
            console.debug("render - no selected symbol");
            return null;
        }


        if (this.state.tickers == null) {
            console.debug("render - no tickers");
            return null;
        }

        console.debug("render - " + this.state.selectedSymbol);

        let tickers = this.state.tickers;
        let optionItems = tickers.map((ticker) =>
            <option value={ticker.symbol}>{ticker.symbol} - {ticker.company}</option>
        );

        let content = (<div className="container-columns-center margin-all">
            <div className="container-rows-left">
                <div className="container-columns-left margin-right">
                    <div className="margin-bottom">
                        <select onChange={this.handleChange}>
                            {optionItems}
                        </select>
                    </div>
                    <SelectSeries selectedSymbol={this.state.selectedSymbol}/>
                </div>
                <div className="container-columns-left wrap margin-left">
                    <div className="margin-bottom chart-height">
                        <GraphStockPrice selectedSymbol={this.state.selectedSymbol}
                                         selectedCompany={this.state.selectedCompany}/>
                    </div>
                    <div className="margin-top chart-height">
                        <GraphPatternDensity selectedSymbol={this.state.selectedSymbol}
                                             selectedCompany={this.state.selectedCompany}/>
                    </div>
                </div>
            </div>
        </div>);

        if (isBrowser) {
            return (
                <BrowserView>{content}</BrowserView>
            );
        }

        return (
            <MobileView>{content}</MobileView>
        );
    }
}

export default SelectTicker;