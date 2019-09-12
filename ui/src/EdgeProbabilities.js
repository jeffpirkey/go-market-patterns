import React from 'react';

class EdgeProbabilities extends React.Component {

    componentDidMount() {
        let url = 'http://localhost:8081/api/latest/edge-probabilities/Low';
        fetch(url).then(response => {
            if (response.status >= 400) {
                throw new Error("Bad response from server");
            }
            return response.json();
        }).then(data => {
            this.setState({probs: data})
        });
    }

    render() {
        if (this.state === null) {
            console.debug("render - no state");
            return null;
        }

        console.debug("render - Low");

        if (this.state.probs === null) {
            return (
                <div className="container-columns-left">
                    <div className="font-title-small">Edge Probabilities</div>
                    <div className="container-columns-left border font-default-text padding-sides">
                        No Edge Probabilities Available
                    </div>
                </div>
            )
        }

        // High
        let bestUpHigh = this.state.probs.bestUpHigh;
        let bestUpHighPerc = ((bestUpHigh.upCount / bestUpHigh.totalCount) * 100).toFixed(2);
        let bestDownHigh = this.state.probs.bestDownHigh;
        let bestDownHighPerc = ((bestDownHigh.upCount / bestDownHigh.totalCount) * 100).toFixed(2);
        let bestNoChangeHigh = this.state.probs.bestNoChangeHigh;
        let bestNoChangeHighPerc = ((bestNoChangeHigh.upCount / bestNoChangeHigh.totalCount) * 100).toFixed(2);
        // Low
        let bestUpLow = this.state.probs.bestUpLow;
        let bestUpLowPerc = ((bestUpLow.upCount / bestUpLow.totalCount) * 100).toFixed(2);
        let bestDownLow = this.state.probs.bestDownLow;
        let bestDownLowPerc = ((bestDownLow.upCount / bestDownLow.totalCount) * 100).toFixed(2);
        let bestNoChangeLow = this.state.probs.bestNoChangeLow;
        let bestNoChangeLowPerc = ((bestNoChangeLow.upCount / bestNoChangeLow.totalCount) * 100).toFixed(2);
        return (
            <div className="container-columns-left">
                <div className="font-title-small">Edge Probabilities</div>
                <div className="container-columns-left border padding-sides">
                    <div className="font-default-text underline">Highest Next Day <b>Up</b></div>
                    <div className="font-text-small">[{bestUpHigh.symbol}] {bestUpHigh.value} {bestUpHighPerc}%</div>
                    <div className="font-default-text underline">Highest Next Day <b>Down</b></div>
                    <div className="font-text-small">[{bestDownHigh.symbol}] {bestDownHigh.value} {bestDownHighPerc}%
                    </div>
                    <div className="font-default-text underline">Highest Next Day <b>No Change</b></div>
                    <div
                        className="font-text-small">[{bestNoChangeHigh.symbol}] {bestNoChangeHigh.value} {bestNoChangeHighPerc}%
                    </div>
                    <div className="font-default-text underline">Lowest Next Day <b>Up</b></div>
                    <div className="font-text-small">[{bestUpLow.symbol}] {bestUpLow.value} {bestUpLowPerc}%</div>
                    <div className="font-default-text underline">Lowest Next Day <b>Down</b></div>
                    <div className="font-text-small">[{bestDownLow.symbol}] {bestDownLow.value} {bestDownLowPerc}%</div>
                    <div className="font-default-text underline">Lowest Next Day <b>No Change</b></div>
                    <div
                        className="font-text-small">[{bestNoChangeLow.symbol}] {bestNoChangeLow.value} {bestNoChangeLowPerc}%
                    </div>
                </div>
            </div>
        );
    }
}

export default EdgeProbabilities;