import update from 'immutability-helper'

export function updateSelectedTicker(state, data) {
    let newState = update(state, {tickerData: {selectedTicker: {$set: data}}});
    return newState;
}

export function updateTickerData(state, data) {
    let newState = update(state, {tickerData: {$set: data}});
    return newState;
}

export function updateTickers(state, data) {
    let newState = update(state, {tickerData: {tickers: {$set: data}, selectedTicker: {$set: data[0]}}});
    return newState;
}

export function updatePrediction(state, data) {
    let newState = update(state, {tickerData: {prediction: {$set: data}}});
    return newState;
}

export function fetchPrediction (container, state) {
    let predictId = state.tickerData.selectedTicker;
    if (predictId !== "") {
        let url = 'http://localhost:7666/api/latest/predict/' + predictId;
        fetch(url).then(response => {
            if (response.status >= 400) {
                throw new Error("Bad response from server");
            }
            return response.json();
        }).then(data => {
                let newState = updateSelectedTicker(state, predictId);
                newState = updatePrediction(newState, data);
            container.setState(newState);
            }
        );
    } else {
        container.setState(state);
    }
}
