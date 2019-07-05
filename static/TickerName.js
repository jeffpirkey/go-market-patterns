import React from 'react';

class TickerName extends React.Component {
    constructor() {
        super();
    }

    render() {
        let tickerNames = this.props.state.tickers;
        let optionItems = tickerNames.map((ticker) =>
            "<option key={tickers.name}>{tickers.name}</option>"
        );

        return (
            "<div><select>{optionItems} </select></div>"
        )
    }
}

export default TickerName;