import React from 'react';
import ReportPrediction from "./ReportPrediction";
import {BrowserView, isBrowser, MobileView} from "react-device-detect";
import GraphPatternDensity from "./GraphPatternDensity";

class SelectSeries extends React.Component {

    constructor(props) {
        super(props);
        this.handleChange = this.handleChange.bind(this);
    }

    componentDidMount() {

        let symbol = this.props.selectedSymbol;
        fetch('http://localhost:8081/api/latest/series/' + symbol)
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            let firstSeriesName = data.series[0].name;
            let firstSeriesLength = data.series[0].length;
            this.setState({
                series: data.series,
                selectedSeries: firstSeriesName,
                selectedLength: firstSeriesLength
            });
        });
    }

    componentDidUpdate(prevProps, prevState, snapshot) {

        if (prevProps.selectedSymbol !== this.props.selectedSymbol) {
            let symbol = this.props.selectedSymbol;
            fetch('http://localhost:8081/api/latest/series/' + symbol)
                .then(response => {
                    if (response.status >= 400) {
                        throw new Error("Bad response from server");
                    }
                    return response.json();
                }).then(data => {
                let firstSeriesName = data.series[0].name;
                let firstSeriesLength = data.series[0].length;
                this.setState({
                    series: data.series,
                    selectedSeries: firstSeriesName,
                    seriesLength: firstSeriesLength
                });
            });
        }
    }

    handleChange(event) {
        let series = this.state.series;
        let selectedSeries = event.target.value;
        let selectedLength = series.filter(s => s.name === selectedSeries)[0].length;
        this.setState({
            series: series,
            selectedSeries: selectedSeries,
            selectedLength: selectedLength
        });
    }

    render() {

        if (this.state === null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.series === null) {
            console.debug("render - no series");
            return null;
        }

        if (this.state.selectedSeries === null || this.state.selectedLength === null) {
            console.debug("render - no selected series");
            return null;
        }

        console.debug("render - " + this.state.selectedSeries);

        let series = this.state.series;
        let optionItems = series.map(s =>
            <option value={s.name}>{s.name}</option>
        );

        let content = (
            <div className="container-rows-left">
                <div className="container-column-left">
                    <div className="margin-bottom">
                        <select onChange={this.handleChange}>
                            {optionItems}
                        </select>
                    </div>
                    <ReportPrediction selectedSymbol={this.props.selectedSymbol}
                                      selectedSeries={this.state.selectedSeries}
                                      selectedLength={this.state.selectedLength}/>
                </div>
                <div className="container-columns-left margin-left">
                    <div className="margin-top chart-height">
                        <GraphPatternDensity selectedSymbol={this.props.selectedSymbol}
                                             selectedCompany={this.props.selectedCompany}
                                             selectedLength={this.state.selectedLength}/>
                    </div>
                </div>
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

export default SelectSeries;