import React, {Component} from 'react';
import ReactDOM from 'react-dom';
import TickerName from './TickerName';

class TickerNameSearch extends Component {
    constructor(props) {
        super(props);
        this.state = {
            tickers: [],
        };
    }
}

componentDidMount()
{
    let initialTickers = [];
    fetch('localhost:7666/api/ticker-names/')
        .then(response => {
            return response.json();
        }).then(data => {
        initialTickers = data.results.map((ticker) => {
            return ticker
        });
        console.log(initialTickers);
        this.setState({
            names: initialTickers,
        });
    });
}

render()
{
    return (
        <TickerName state={this.state}/>
    );
}

// after component is finished

export default TickerNameSearch;

ReactDOM.render(<TickerNameSearch/>, document.getElementById('react-search'));