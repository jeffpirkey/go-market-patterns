import React from 'react';

class PredictResult extends React.Component {

    constructor(props) {
        super(props);
    }

    render() {

        if (this.props.tickerData.prediction == null) {
            console.debug ("empty rendering")
            return null;
        }

        console.debug ("rendering " + this.props.tickerData.selectedTicker)

        let ticker = this.props.tickerData.prediction.ticker;
        let fromDate = this.props.tickerData.prediction.predictingFromDate;
        let nextDate = this.props.tickerData.prediction.predictingDate;
        let seriesAry = this.props.tickerData.prediction.series;

        let probItems = seriesAry.map(function (data) {

                let probAry = [];
                let probMap = data.probabilityOfNextBeing;
                Object.keys(probMap).forEach(function (key) {
                    let perc = Number.parseFloat(probMap[key] * 100).toFixed(2);
                    probAry.push(<div>{key} = {perc}%</div>)
                });

                return (<div class="App-container-left">
                    Predictions
                    <div class="App-border-container-left">
                        <div class="App-small-title">{data.name}</div>
                        <div class="App-small-content">Prior periods were {data.priorPeriodsWere}</div>
                        <div class="App-small-content">Probabilities:</div>
                        <div class="App-small-content">{probAry}</div>
                    </div>
                </div>)
            }
        );


        return (
            <div class="App-container-left">
                <div class="App-title">{ticker}</div>
                <div class="App-content">From Date: {fromDate}</div>
                <div class="App-content">Prediction Date: {nextDate}</div>
                <div class="App-content">{probItems}</div>
            </div>
        )
    }
}

export default PredictResult;