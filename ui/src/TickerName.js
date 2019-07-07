import React from 'react';

class TickerName extends React.Component {
    constructor(props) {
        super(props);

        this.handleChange = this.handleChange.bind(this);
        this.handleSubmit = this.handleSubmit.bind(this);
    }

    handleChange(event) {
        this.setState({predict: event.target.value});
    }

    handleSubmit(event) {
        let initialTickers = [];

        fetch('http://localhost:7666/api/latest/predict/' + this.state.predict)
            .then(response => {
                if (response.status >= 400) {
                    throw new Error("Bad response from server");
                }
                return response.json();
            }).then(data => {
            this.setState({
                prediction: data
            });
        });

        event.preventDefault();
    }

    render() {
        let tickers = this.props.state.tickers;
        let optionItems = tickers.map((ticker) =>
            <option key={ticker}>{ticker}</option>
        );

        return (
            <form onSubmit={this.handleSubmit}>
                <label className="App-label">Ticker Name</label>
                <select value={this.props.state.value} onChange={this.handleChange}>
                    {optionItems}
                </select>

                <input type="submit" value="Submit"/>
            </form>
        )
    }
}

export default TickerName;