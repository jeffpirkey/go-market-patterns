import React from 'react';
import HighchartsReact from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import darkUnica from 'highcharts/themes/dark-unica'

darkUnica(Highcharts);

class GraphPatternDensity extends React.Component {

    componentDidMount() {
        // Initial setup
        let symbol = this.props.selectedSymbol;
        let selectedLength = this.props.selectedLength;
        if (symbol && selectedLength) {
            let url = 'http://localhost:8081/api/latest/graph/pattern-density/' +
                symbol + '?length=' + selectedLength;
            fetch(url).then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
                this.setState({
                    selectedSymbol: data.symbol,
                    selectedCompany: data.companyName,
                    selectedLength: selectedLength,
                    data: data.graphData
                });
            });
        } else {
            console.error("no symbol or length")
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {

        if (prevProps.selectedSymbol !== this.props.selectedSymbol ||
            prevProps.selectedLength !== this.props.selectedLength) {
            let symbol = this.props.selectedSymbol;
            let selectedLength = this.props.selectedLength;
            if (symbol && selectedLength) {
                let url = 'http://localhost:8081/api/latest/graph/pattern-density/' +
                    symbol + '?length=' + selectedLength;
                fetch(url).then(response => {
                    if (response.status >= 400) {
                        throw new Error("Bad response from server");
                    }
                    return response.json();
                }).then(data => {
                    this.setState({
                        selectedSymbol: data.symbol,
                        selectedCompany: data.companyName,
                        selectedLength: selectedLength,
                        data: data.graphData
                    });
                });
            } else {
                console.error("no symbol or length")
            }
        }
    }

    render() {

        if (this.state === null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.data === null) {
            console.debug("render - no data");
            return null;
        }

        if (this.props.selectedSymbol === null) {
            console.debug("render - no symbol");
            return null;
        }

        if (this.props.selectedCompany === null) {
            console.debug("render - no company");
            return null;
        }

        if (this.props.selectedLength === null) {
            console.debug("render - no series length");
            return null;
        }

        console.debug("render - " + this.props.selectedSymbol + " " + this.props.selectedCompany +
            " " + this.props.selectedLength);

        let title = "(" + this.props.selectedSymbol + ") " + this.props.selectedCompany + " Pattern " +
            this.props.selectedLength + " Period Density";
        let chartOptions = {
            title: {text: title},
            chart: {height: 300, type: 'column'},
            plotOptions: {
                column: {
                    stacking: 'normal',
                    dataLabels: {
                        enabled: true,
                        color: (Highcharts.theme && Highcharts.theme.dataLabelsColor) || 'white'
                    }
                }
            },
            xAxis: {categories: this.state.data.categories},
            series: [{name: 'Totals', data: this.state.data.totals},
                {name: 'Up', data: this.state.data.ups},
                {name: 'Down', data: this.state.data.downs},
                {name: 'No Change', data: this.state.data.nochanges}],
        };

        return (
            <HighchartsReact
                highcharts={Highcharts}
                constructorType={'chart'}
                options={chartOptions}/>
        );
    }
}

export default GraphPatternDensity