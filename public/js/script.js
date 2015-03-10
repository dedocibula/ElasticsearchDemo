(function($, Mustache, window, document, undefined) {
	var Monitor = (function() {
		
		function Monitor(config, templates) {
			this.socket = null;
			this.initialized = false;

			this.config = $.extend({}, config);
			this.templates = $.extend({}, templates);
		}

		Monitor.prototype = {
			initialize: function() {
				var self = this;

				Mustache.tags = self.config.mustache.tags;
				self._setupCallbacks();
				self._parseTemplates();
				self._setupWebSocket();

				self.initialized = true;
			},

			fetchResults: function() {
				var self = this;

				if (self.initialized) {
					self._fetch().done(function(data) {
				  		var resultsTemplate = self.templates.results;

			  			var rendered = Mustache.render(resultsTemplate.body, self._formResults(data));
			  			var container = self._extractContainer(resultsTemplate.container);
			  			container.html(rendered);

			  			self._afterFetch(data.length > 0);
		  			});
				} else {
					throw 'Monitor not initialized';
				}
			},

			_setupCallbacks: function() {
				var self = this;

				var clearBtn = $('#clearBtn');
				clearBtn.on('click', function() {
			  		clearBtn.attr('disabled', 'disabled');
			  	});
			  	self._enableClearButton = function() {
			  		clearBtn.removeAttr('disabled');
			  	};
			},

			_parseTemplates: function() {
				var self = this;

				$.each(self.templates, function(id, template) {
					var body = template.body.html();
					Mustache.parse(body);
					self.templates[id].body = body;
				});
			},

			_setupWebSocket: function() {
				var self = this;

				var socket = new WebSocket(self.config.endpoints.ws);

				socket.onmessage = function(message) {
			  		var data = JSON.parse(message.data);
			  		var attemptTemplate = self.templates.attempt;

			  		var rendered = Mustache.render(attemptTemplate.body, data.Attempt);
			  		var container = self._extractContainer(attemptTemplate.container, data.Success);
			  		container.prepend(rendered);

			  		self._truncate(container.children(), self.config.maxAttemptCount);

			  		self._afterMessage(data.Success);
			  	};

			  	self.socket = socket;
			},

			_fetch: function() {
				return $.ajax({
					url: this.config.endpoints.results,
					dataType: 'json'
				});
			},

			_formResults: function(data) {
				return {
			  		'results': $.map(data, function(doc, index) {
			  			doc['Index'] = index + 1;
			  			return doc;
			  		}), 
			  		'if': function() {
			  			return function(text, render) {
			  				return data.length > 0 ? render(text) : '';
			  			}
			  		} 
			  	};
			},

			_extractContainer: function(containerValue, success) {
				if (containerValue instanceof $)
					return containerValue;
				else if (typeof containerValue === 'function' && typeof success === 'boolean')
					return containerValue(success);
				else
					throw 'Unrecognized input';
			},

			_truncate: function(children, maxCount) {
				if (children.length > maxCount) {
					$(children).last().remove();
				}
			},

			_afterFetch: function(ok) {
				if (ok) self._enableClearButton();
			},

			_afterMessage: function(ok) {
				if (ok) self.fetchResults();
			}
		};

		return Monitor;
	}());

	var config = {
		mustache: {
			tags: [ '[[', ']]' ]
		},
		maxAttemptCount: 10,
		endpoints: {
			ws: 'ws://' + window.location.host + '/Admin/AttemptEndpoint',
			results: '/Admin/ResultEndpoint'
		}
	}

	var templates = {
		results: {
			body: $('#results-template'),
			container: $('#results')
		},
		attempt: {
	  		body: $('#attempt-template'),
	  		container: function(success) {
	  			return success ? $('#successes') : $('#failures');
	  		}
	  	}
	};


  	var monitor = new Monitor(config, templates);
  	monitor.initialize();
  	monitor.fetchResults();
  	
})(jQuery, Mustache, window, document);