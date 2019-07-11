import React from 'react';
import PredictResult from "./PredictResult";
import * as utils from './tickerData'

class TickerName extends React.Component {

    constructor(props) {
        super(props);
        this.state = {
            tickerData: {
                selectedTicker: "",
                tickers: [],
                prediction: null
            }
        };
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount() {

        fetch('http://localhost:7666/api/latest/ticker-names/')
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            let state = utils.updateTickers(this.state, data.names);
            utils.fetchPrediction(this, state)
        });
    }

    handleChange(event) {
        let predictId = event.target.value;
        if (predictId === "") {
            console.error("Invalid predict id: " + predictId);
            return
        }

        let url = 'http://localhost:7666/api/latest/predict/' + predictId;
        fetch(url).then(response => {
            if (response.status >= 400) {
                throw new Error("Bad response from server");
            }
            return response.json();
        }).then(data => {
                let newState = utils.updateSelectedTicker(this.state, predictId);
                newState = utils.updatePrediction(newState, data);
                this.setState(newState);
            }
        );
    }

    render() {
        console.debug("rendering " + this.state.tickerData.selectedTicker)
        let tickers = this.state.tickerData.tickers;
        let optionItems = tickers.map((ticker) =>
            <option key={ticker}>{ticker}</option>
        );

        return (
            <div class="App-container-center">
                <div class="App-container-left">
                    <div class="App-content">Select Ticker:</div>
                    <select value={this.state.tickerData.selectedTicker} onChange={this.handleChange}>
                        {optionItems}
                    </select>
                </div>
                <PredictResult tickerData={this.state.tickerData}/>
            </div>
        )
    }
}

export default TickerName;