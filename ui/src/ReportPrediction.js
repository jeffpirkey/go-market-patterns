import React from 'react';

class ReportPrediction extends React.Component {

    componentDidMount() {

        let predictId = this.props.selectedSymbol;
        if (predictId) {
            let url = 'http://localhost:8081/api/latest/predict/' + predictId;
            fetch(url).then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
                this.setState({prediction: data})
            });
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {

        if (prevProps.selectedSymbol !== this.props.selectedSymbol) {
            let predictId = this.props.selectedSymbol;
            if (predictId) {
                let url = 'http://localhost:8081/api/latest/predict/' + predictId;
                fetch(url).then(response => {
                    if (response.status >= 400) {
                        throw new Error("Bad response from server");
                    }
                    return response.json();
                }).then(data => {
                    this.setState({prediction: data})
                });
            }
        }
    }

    render() {

        if (this.state == null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.prediction == null) {
            console.debug("render - no prediction");
            return null;
        }

        console.debug("render - " + this.props.selectedSymbol);

        let fromDate = this.state.prediction.predictingFromDate;
        let nextDate = this.state.prediction.predictingDate;
        let seriesAry = this.state.prediction.series;
        if (seriesAry == null) {
            return (
                <div className="container-columns-left">
                    <div className="font-title-small">Predictions</div>
                    <div className="container-columns-left border font-default-text padding-sides">No Predictions Available</div>
                </div>
            )
        }

        let probItems = seriesAry.map(function (data) {
                let priorPattern = data.priorPeriodsWere;
                let name = data.name;
                let up = data.probabilityOfNextBeingUp * 100;
                let probUp = "Up = " + up.toFixed(2) + "%";
                let down = data.probabilityOfNextBeingDown * 100;
                let probDown = "Down = " + down.toFixed(2) + "%";
                let nc = data.probabilityOfNextBeingNoChange * 100;
                let probNoChange = "No Change = " + nc.toFixed(2) + "%";
                return (
                    <div className="container-columns-left">
                        <div className="font-title-small">Predictions</div>
                        <div className="container-columns-left border padding-sides">
                            <div className="font-text-small">{name}</div>
                            <div className="font-text-small">{fromDate} was {priorPattern}</div>
                            <div className="font-text-small">{nextDate} Close Probability</div>
                            <div className="font-text-small">{probUp}</div>
                            <div className="font-text-small">{probDown}</div>
                            <div className="font-text-small">{probNoChange}</div>
                        </div>
                    </div>
                )
            }
        );

        return (<div>{probItems}</div>);
    }
}

export default ReportPrediction;