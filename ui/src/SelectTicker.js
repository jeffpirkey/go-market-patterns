import React from 'react';
import SelectSeries from "./SelectSeries";
import {BrowserView, isBrowser, MobileView} from "react-device-detect";
import GraphStockPrice from "./GraphStockPrice";
import EdgeProbabilities from "./EdgeProbabilities";

class SelectTicker extends React.Component {

    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount() {

        fetch('http://localhost:8081/api/latest/tickers/')
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            let firstSymbol = data.tickers[0].symbol;
            let firstCompany = data.tickers[0].company;
            this.setState(
                {
                    selectedSymbol: firstSymbol,
                    selectedCompany: firstCompany,
                    tickers: data.tickers
                }
            );
        });
    }

    handleChange(event) {
        let tickers = this.state.tickers;
        let companyName = tickers.filter(ticker => ticker.symbol === event.target.value)[0].company;
        this.setState(
            {
                selectedSymbol: event.target.value,
                selectedCompany: companyName,
                tickers: tickers
            }
        );
    }

    render() {
        if (this.state === null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.selectedSymbol === null) {
            console.debug("render - no selected symbol");
            return null;
        }

        if (this.state.tickers === null) {
            console.debug("render - no tickers");
            return null;
        }

        console.debug("render - " + this.state.selectedSymbol);

        let tickers = this.state.tickers;
        let optionItems = tickers.map(ticker =>
            <option value={ticker.symbol}>{ticker.symbol} - {ticker.company}</option>
        );

        let content = (
            <div className="container-columns-center margin-all">
                <div className="container-rows-left">
                    <div className="container-columns-left margin-right">
                        <div className="margin-bottom">
                            <select onChange={this.handleChange}>
                                {optionItems}
                            </select>
                        </div>
                    </div>
                    <div className="margin-bottom chart-height">
                        <GraphStockPrice selectedSymbol={this.state.selectedSymbol}
                                         selectedCompany={this.state.selectedCompany}/>
                    </div>
                    <div className="margin-bottom margin-left">
                        <EdgeProbabilities/>
                    </div>
                </div>
                <SelectSeries selectedSymbol={this.state.selectedSymbol}
                              selectedCompany={this.state.selectedCompany}/>
            </div>
        );

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