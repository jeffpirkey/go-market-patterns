import React, {Component} from 'react';
import TickerName from './TickerName';
import PredictResult from "./PredictResult";

class TickerNameSearch extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            tickers: [],
            prediction: {}
        };
    }

    componentDidMount() {
        let initialTickers = [];
        fetch('http://localhost:7666/api/latest/ticker-names/')
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            this.setState({
                tickers: data.names,
                predict: data.names[0]
            });
        });
    }

    render() {
        console.debug("rendering " + this.state)
        return (
            <TickerName state={this.state}/>
            <PredictResult/>
        );
    }
}

// after component is finished

export default TickerNameSearch;
