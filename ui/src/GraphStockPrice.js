import React, {Component} from 'react';
import HighchartsReact from 'highcharts-react-official';
import Highcharts from 'highcharts/highstock';
import darkUnica from 'highcharts/themes/dark-unica'

darkUnica(Highcharts);

class GraphStockPrice extends Component {

    componentDidMount() {
        let predictId = this.props.selectedSymbol;
        if (predictId) {
            let url = 'http://localhost:8081/api/latest/graph/stock/' + predictId;
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
                let url = 'http://localhost:8081/api/latest/graph/stock/' + predictId;
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

        let title = "(" + this.props.selectedSymbol + ") " + this.props.selectedCompany + " NYSE";
        let chartOptions = {
            chart: {height: 300},
            title: {text: title},
            series: this.state.data,
        };

        return (
            <HighchartsReact
                highcharts={Highcharts}
                constructorType={'stockChart'}
                options={chartOptions}
            />
        );
    }
}

export default GraphStockPrice