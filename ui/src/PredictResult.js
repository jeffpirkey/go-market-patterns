import React from 'react';

class PredictResult extends Component {
    constructor(props) {
        super(props);
    }

    render () {
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

export default PredictResult;