import React from 'react';
import ReportPrediction from "./ReportPrediction";
import {BrowserView, isBrowser, MobileView} from "react-device-detect";

class TickerName extends React.Component {

    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount() {

        let predictId = this.props.selectedSymbol;
        fetch('http://localhost:7666/api/latest/series/' + predictId)
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            this.setState({seriesNames: data.names, selectedSeries: data.names[0]});
        });
    }

    componentDidUpdate(prevProps, prevState, snapshot) {

        if (prevProps.selectedSymbol !== this.props.selectedSymbol) {
            let predictId = this.props.selectedSymbol;
            fetch('http://localhost:7666/api/latest/series/' + predictId)
                .then(response => {
                    if (response.status >= 400) {
                        throw new Error("Bad response from server");
                    }
                    return response.json();
                }).then(data => {
                this.setState({seriesNames: data.names, selectedSeries: data.names[0]});
            });
        }
    }

    handleChange(event) {
        this.setState({selectedSymbol: event.target.value});
        console.debug("tmp");
    }

    render() {

        if (this.state == null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.seriesNames == null) {
            console.debug("render - no series");
            return null;
        }

        if (this.state.selectedSeries == null) {
            console.debug("render - no selected series");
            return null;
        }

        console.debug("render - " + this.state.selectedSeries);

        let seriesNames = this.state.seriesNames;
        let optionItems = seriesNames.map((name) =>
            <option value={name}>{name}</option>
        );

        let content = (
            <div className="container-columns-left">
                <div className="margin-bottom">
                    <select onChange={this.handleChange}>
                        {optionItems}
                    </select>
                </div>
                <ReportPrediction selectedSymbol={this.props.selectedSymbol}
                                  selectedSeries={this.state.selectedSeries}/>
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

export default TickerName;