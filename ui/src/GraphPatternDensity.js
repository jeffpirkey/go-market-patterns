import React, {Component} from 'react';
import HighchartsReact from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import darkUnica from 'highcharts/themes/dark-unica'

darkUnica(Highcharts);

class StockPriceGraph extends Component {

    componentDidMount() {
        let predictId = this.props.selectedSymbol;
        if (predictId) {
            let url = 'http://localhost:8081/api/latest/graph/pattern-density/' + predictId;
            fetch(url).then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
                this.setState({data: data})
            });
        }
    }

    componentDidUpdate(prevProps, prevState, snapshot) {

        if (prevProps.selectedSymbol !== this.props.selectedSymbol) {
            let predictId = this.props.selectedSymbol;
            if (predictId) {
                let url = 'http://localhost:8081/api/latest/graph/pattern-density/' + predictId;
                fetch(url).then(response => {
                    if (response.status >= 400) {
                        throw new Error("Bad response from server");
                    }
                    return response.json();
                }).then(data => {
                    this.setState({data: data})
                });
            }
        }
    }

    render() {

        if (this.state == null) {
            console.debug("render - no state");
            return null;
        }

        if (this.state.data == null) {
            console.debug("render - no data");
            return null;
        }

        console.debug("render - " + this.props.selectedSymbol);

        let title = "(" + this.props.selectedSymbol + ") " + this.props.selectedCompany + " Pattern Density";
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
                options={chartOptions}
            />
        );
    }
}

export default StockPriceGraph