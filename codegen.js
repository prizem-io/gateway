#! /usr/local/bin/node

var fs = require('fs');
var path = require('path');
var _ = require('lodash');
var inflection = require('inflection');

var nodegen = require("oas-nodegen");
var utils = nodegen.utils;

var loader = nodegen.createLoader();
var go = require('./go');
var modules = nodegen.createModules()
  .registerModule(go);

var apiSpecLocations = _
  .chain(fs.readdirSync(path.resolve(__dirname, 'specs')))
  .filter(function(name) {
    return name != 'common.yaml';
  })
  .map(function(name) {
    return path.resolve(__dirname, 'specs', name);
  })
  .value();

var config = {
  packages : {
    model : 'github.com/prizem-io/gateway/models',
    resource : 'github.com/prizem-io/gateway/apis',
    service : 'github.com/prizem-io/gateway/services'
  }
}

// supported validation rules from https://github.com/asaskevich/govalidator
var validators = [
  "alpha", "alphanum", "ascii", "base64", "creditcard", "datauri",
  "dialstring", "dns", "email", "float", "fullwidth", "halfwidth",
  "hexadecimal", "hexcolor", "host", "int", "ip", "ipv4", "ipv6",
  "isbn10", "isbn13", "json", "latitude", "longitude", "lowercase",
  "mac", "multibyte", "null", "numeric", "port", "printableascii",
  "requri", "requrl", "rgbcolor", "ssn", "semver", "uppercase", "url",
  "utfdigit", "utfletter", "utfletternum", "utfnumeric", "uuid",
  "uuidv3", "uuidv4", "uuidv5", "variablewidth"
]

var writer = nodegen.createWriter(__dirname);

var templates = nodegen.createTemplates()
  .setDefaultOptions({ noEscape : true })
  .registerTemplateDirectory(path.resolve(__dirname, 'templates'));

function nameSortFn(a, b) {
  if (a.name < b.name) {
    return -1;
  } else if (a.name > b.name) {
    return 1;
  } else {
    return 0;
  }
};

function setVarNames(obj) {
  obj.displayName = inflection.humanize(obj.name, true);
  obj.singularVar = inflection.singularize(utils.uncapitalize(obj.name));
  obj.pluralVar = inflection.pluralize(utils.uncapitalize(obj.name));
}

///////////////////////////////////////////////////////////
// MODEL GENERATOR
///////////////////////////////////////////////////////////

var additionalFields = {
  PluginConfig: [
    'Config interface{}'
  ],
  Operation: [
    'UpstreamConfig interface{}'
  ]
}

var modelGenerator = nodegen.createGenerator()
  .configure({
    modelPackage : config.packages.model
  })
  .use(modules.get('Go'));

modelGenerator.setTypeFormat = "[]{type}";

var modelTemplate = templates.compileFromFile('models.handlebars');

modelGenerator.onDecorate('Model', function(context) {
  var model = context.model;

  if (_.contains(['PaginatedList', 'Entity', 'Auditable'], model.name)) {
    model.abstract = true;
  }

  model.tableName = inflection.underscore(model.name, true);
  if (model.tableName.endsWith('_update')) {
    model.update = true;
    model.tableName = model.tableName.slice(0, -7);
  }

  var entity = _.find(model.references, function(i) {
    return _.includes(['Entity'], i.name);
  });
  if (entity != null) {
    model.entity = entity.name;
  }

  setVarNames(model);
  model.additionalFields = additionalFields[model.name];
});

modelGenerator.onDecorate('Property', function(context) {
  var model = context.model;
  var property = context.property;

  this.addStructTag('json', property.name, property);
  this.addStructTag('yaml', property.name, property);
  this.addStructTag('msgpack', property.name, property);

  var rules = []
  if (property.required) {
    rules.push('required');
  }
  if (property.type == 'string') {
    if (property.maxLength) {
      rules.push('stringlength(' + (property.minLength || 0) + '|' + property.maxLength + ')');
    } else if (property.minLength) {
      rules.push('stringlength(' + property.minLength + '|1024)');
    }
    if (property.pattern) {
      rules.push('matches(' + property.pattern + ')');
    }
  }
  if (property['x-validate']) {
    var rs = [];
    for (var i=0; i<property['x-validate'].length; i++) {
      var v = property['x-validate'][i];
      if (validators.indexOf(v) >= 0) {
        rs.push(v);
      } else {
        console.log("Unrecognized validation rule for " + model.name + "." + columnName + ": " + v)
      }
    }
    rules = rules.concat(rs);
  }
  if (rules.length > 0) {
    this.addStructTag('valid', rules.join(","), property);
  }
});

modelGenerator.onWrite('Context', function(context, name) {
  var data = {
    package: this.config.modelPackage,
    imports: _.union.apply(null,
      _.map(context.models, function(model) { return model.imports; }))
      .sort(),
    models: _.filter(context.models, function(model) {
      return !model.name.endsWith('List');
    }).sort(nameSortFn),
    enums: utils.sortKeys(context.enums)
  }

  var code = modelTemplate(data);
  writer.write('generated.go', code);
});

///////////////////////////////////////////////////////////

loader.load(apiSpecLocations).then(function(specifications) {
  try {
    modelGenerator.process(specifications);
  } catch (error) {
    console.log("Processing failure", error.stack);
    process.exit(1);
  }
}).fail(function(error) {
  console.log("Loading failure", error.stack);
  process.exit(1);
});
