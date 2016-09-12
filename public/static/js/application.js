class NavigationBar extends React.Component {
  constructor() {
    super();
  }

  render() {
    return (
      <nav className="navbar navbar-inverse navbar-static-top">
        <div className="container">
          <div className="navbar-header">
            <a className="navbar-brand" href="#">StockExchange</a>
          </div>
          <div id="navbar" className="collapse navbar-collapse">
            <div className="navbar-form navbar-right">
              <div className="input-group">
                <input type="text" ref="searchBox" className="form-control" placeholder="Search for..."/>
                <span className="input-group-btn">
                  <button className="btn btn-default" type="button" onClick={() => this.props.onSearchClick(this.refs.searchBox.value)}>
                    <span className="glyphicon glyphicon-search" aria-hidden="true"></span>
                  </button>
                </span>
              </div>
            </div>
          </div>
        </div>
      </nav>
    );
  }
}

class Portfolio extends React.Component {
  constructor() {
    super();
  }

  render() {
    if(!this.props.value.shares || this.props.value.shares.length == 0) {
      return (
        <div className="panel panel-info">
          <div className="panel-heading">
            <h3 className="panel-title">
              <span>Portfolio</span>
            </h3>
          </div>
          <div className="panel-body">
            <h4 className="text-center">Your portfolio is empty</h4>
          </div>
          <div className="panel-footer">
            <span className="text-right text-uppercase"><strong>Cash ${this.props.value.balance}</strong></span>
          </div>
        </div>
      );
    }

    return (
      <div className="panel panel-info">
        <div className="panel-heading">
          <h3 className="panel-title">
            <span className="text-uppercase">Portfolio</span>
          </h3>
        </div>
        <div className="panel-body">
          <table className="table table-hover">
            <thead>
              <tr>
                <th>Symbol</th>
                <th>Name</th>
                <th>Paid Price</th>
                <th>Quantity</th>
              </tr>
            </thead>
            <tbody>
              {
                this.props.value.shares.map(function(item) {
                  return (
                    <tr key={item.symbol}>
                      <th scope="row">{item.symbol}</th>
                      <td>{item.name}</td>
                      <td>{item.paidPrice}</td>
                      <td>{item.quantity}</td>
                    </tr>
                  );
                })
              }
            </tbody>
          </table>
        </div>
        <div className="panel-footer">
          <span className="text-right text-uppercase"><strong>Cash ${this.props.value.balance}</strong></span>
        </div>
      </div>
    );
  }
}

class SearchContainer extends React.Component {
  constructor() {
    super();
    this.validate = this.validate.bind(this);
  }

  validate(e) {
    let charCode = (e.which) ? e.which : event.keyCode

    if (charCode > 31 && (charCode < 48 || charCode > 57)) {
      e.preventDefault();
    }
  }

  render() {
    if(this.props.message) {
      return (
        <div className="panel panel-danger">
          <div className="panel-body">
            <h4 className="text-center">{this.props.message}</h4>
          </div>
        </div>
      );
    }

    if(this.props.result && this.props.result.length > 0) {
      return (
        <div className="panel panel-warning">
          <div className="panel-heading">
            <h3 className="panel-title">
              <span>Stock</span>
            </h3>
          </div>
          <div className="panel-body">
            <table className="table table-hover">
              <thead>
                <tr>
                  <th>Symbol</th>
                  <th>Name</th>
                  <th>Ask Price</th>
                  <th>Bid Price</th>
                </tr>
              </thead>
              <tbody>
                {
                  this.props.result.map(function(item) {
                    return (
                      <tr key={item.symbol}>
                        <th scope="row">{item.symbol}</th>
                        <td>{item.name}</td>
                        <td>${item.askPrice}</td>
                        <td>${item.bidPrice}</td>
                      </tr>
                    );
                  })
                }
              </tbody>
            </table>
          </div>
          <div className="panel-footer">
            <div className="input-group">
              <div className="input-group">
                <span className="input-group-addon">Quantity</span>
                <input type="number" pattern="[0-9]" min="1" defaultValue="1" required="required" className="form-control" aria-label="hidden" ref="quantityBox" onKeyPress={this.validate} />
                <div className="input-group-btn">
                  <button type="button" className="btn btn-default" onClick={() => this.props.onBuyClick(this.props.result[0], this.refs.quantityBox.value)}>Buy</button>
                  <button type="button" className="btn btn-default" onClick={() => this.props.onSellClick(this.props.result[0], this.refs.quantityBox.value)}>Sell</button>
                </div>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return <div></div>;
  }
}

class Alert extends React.Component {
  constructor() {
    super()
    this.close = this.close.bind(this);
    this.state = {visible: true};
  }

  render() {
    if(this.state.visible && this.props.message) {
      return (
        <div className={"alert alert-master alert-" + this.props.level} role="alert">
          <button type="button" className="close" onClick={this.close}>×</button>
          <span>{this.props.message}</span>
        </div>
      );
    }
    return <div></div>;
  }

  close() {
    this.setState({visible: false});
  }
}

class ProgressBar extends React.Component {
  constructor() {
    super();
  }

  render() {
    if(!this.props.visible) {
      return <div></div>;
    }

    var style = {width: "100%"};
    return (
      <div className="progress">
        <div style={style} className="progress-bar progress-bar-info progress-bar-striped active" role="progressbar" aria-valuenow="40" aria-valuemin="0" aria-valuemax="100">
          <span className="sr-only">40% Complete (success)</span>
        </div>
      </div>
    );
  }
}

class ApplicationContainer extends React.Component {
  constructor() {
    super();

    this.state = {
      waiting: false,
      portfolio: {balance: 100000, shares:[]},
      search:{result: [], message: ""},
      alert: {level: "info", message: ""},
    };

    this.search = this.search.bind(this);
    this.buy = this.buy.bind(this);
    this.sell = this.sell.bind(this);
  }

  render() {
    return (
      <div>
        <ProgressBar visible={this.state.waiting} />
        <NavigationBar onSearchClick={this.search} />
        <div className="container container-small">
          <Portfolio value={this.state.portfolio} />
          <SearchContainer result={this.state.search.result} message={this.state.search.message} onBuyClick={this.buy} onSellClick={this.sell}/>
          <Alert level={this.state.alert.level} message={this.state.alert.message}/>
        </div>
      </div>
    );
  }

  componentDidMount() {
    $.ajax({
      url: "/api/v1/portfolio",
      dataType: 'json',
      error: function(jqXHR, textStatus, errorThrown) {
        this.setState({waiting: false, alert: {level: "danger", message: jqXHR.responseText}});
      }.bind(this),
      success: function(portfolio) {
        this.setState({waiting: false, portfolio: portfolio});
      }.bind(this)
    });
  }

  search(text) {
    $.ajax({
      url: "/api/v1/search?query="+text,
      dataType: 'json',
      beforeSend: function() {
        this.setState({waiting: true})
      }.bind(this),
      error: function() {
        this.setState({waiting: false, portfolio: this.state.portfolio, search:{result: [], message: "Nothing has been found"}})
      }.bind(this),
      success: function(stock) {
        this.setState({waiting: false, portfolio: this.state.portfolio, search:{result: stock, message: ""}})
      }.bind(this)
    });
  }

  buy(stock, quantity) {
    $.ajax({
      url: "/api/v1/buy?quantity="+quantity,
      contentType: "application/json; charset=utf-8",
      type: "POST",
      dataType: 'json',
      data:  JSON.stringify(stock),
      beforeSend: function() {
        this.setState({waiting: true})
      }.bind(this),
      error: function(jqXHR, textStatus, errorThrown) {
        this.setState({waiting: false, alert: {level: "danger", message: jqXHR.responseText}});
      }.bind(this),
      success: function(portfolio) {
        this.setState({waiting: false, portfolio: portfolio, alert: {level: "info", message: ""}})
      }.bind(this)
    });
  }

  sell(stock, quantity) {
    let invoice = {
      Symbol: stock.symbol,
      Price: stock.bidPrice,
      Quantity: parseInt(quantity),
    };
    $.ajax({
      url: "/api/v1/sell",
      contentType: "application/json; charset=utf-8",
      type: "POST",
      dataType: 'json',
      data:  JSON.stringify(invoice),
      beforeSend: function() {
        this.setState({waiting: true})
      }.bind(this),
      error: function(jqXHR, textStatus, errorThrown) {
        this.setState({waiting: false, alert: {level: "danger", message: jqXHR.responseText}});
      }.bind(this),
      success: function(portfolio) {
        this.setState({waiting: false, portfolio: portfolio, alert: {level: "info", message: ""}})
      }.bind(this)
    });
  }
}

ReactDOM.render(<ApplicationContainer/>, document.getElementById('app'));
