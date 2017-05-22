/*
 * Copyright 2016 Capital One Services, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

var _ = require('lodash');
var inflection = require('inflection');
var util = require('oas-nodegen').Utilities;//require('../utilities');

module.exports.name = 'Go';

module.exports.dependsOn = ['Helpers'];

module.exports.initialize = function() {
  this.translateTypeMap = {
    // Standard type values
    'integer' : 'int',
    'number' : 'float64',
    'string' : 'string',
    'boolean' : 'bool',
    'File' : 'File',
    // Non-standard values that can be assigned in x-type but not type
    'any' : 'interface{}'
  };

  this.translateFormatMap = {
    'int32' : 'int32',
    'int64' : 'int64',
    'float' : 'float32',
    'double' : 'float64',
    'byte' : 'string',
    'date' : 'time.Time',
    'date-time' : 'time.Time'
  };

  this.imports = [
  ];

  this.convertMap = {
    'File' : 'InputStream'
  };

  this.paramTypeMap = {
    'formData' : 'form'
  };

  this.typePrefixPackages = {
  };

  this.setTypeFormat = "map[{type}]struct{}";

  this.convertEnums = true;
  this.enumSuffix = '';

  this.addImport = function(package, imports) {
    if (!package || !imports) {
      return;
    }

    if (_.isObject(imports) && !_.isArray(imports)) {
      imports = imports.imports = imports.imports || [];
    }

    if (!_.includes(imports, package)) {
      imports.push(package);
    }
  }

  this.addTypeImport = function(type, imports) {
    if (!type || !imports) {
      return;
    }

    if (type.startsWith('*')) {
      type = type.substring(1);
    }

    var index = type.indexOf('.');
    if (index != -1) {
      var package = type.substring(0, index);
      package = this.typePrefixPackages[package] || package;

      this.addImport(package, imports);
    }
  }

  this.getPackagePrefix = function(value) {
    var index = value.lastIndexOf('/');
    if (index != -1) {
      value = value.substring(index + 1);
    }

    return value;
  }

  this.addStructTag = function(tagName, tagValue, tags) {
    if (arguments.length == 2) {
      tags = tagValue;
      tagValue = tagName;
      tagName = '_'
    }
    if (!tagName || !tagValue || !tags) {
      return;
    }

    if (_.isObject(tags) && !_.isArray(tags)) {
      tags = tags.tags = tags.tags || {};
    }

    var dest = tags[tagName];
    if (!dest) {
      dest = [];
      tags[tagName] = dest;
    }

    if (!_.includes(dest, tagValue)) {
      dest.push(tagValue);
    }
  }

  this.joinStrings = function(values) {
    return _.map(values, _.bind(function(v) { return this.escapeGoString(v); }, this)).join(', ');
  }

  this.typeTranslators = [];

  this.addTypeTranslator = function(translator) {
    this.typeTranslators.push(translator);

    return this; // Allow chaining
  }

  this.overrideModelPackage = function($ref, resolved) {
    return null;
  }

  this.customSliceHandler = function(context, schema, itemType) {
    return null;
  }

  this.translateType = function(schema, imports, context, external, noPointer) {
    var spec = context.spec;
    var references = context.references;

    if (schema == null) {
      return null;
    }

    var type = null;

    _.each(this.typeTranslators, _.bind(function(translator) {
      type = translator.call(this, schema, imports);
      return type != null;
    }, this));

    if (type != null) {
      this.addTypeImport(type, imports);
      return type;
    }

    var modelPackage = this.config.modelPackage || this.config.package || 'model';

    var schemaType = schema['x-type'] || schema.type;

    if (schemaType) {
      if (schemaType == 'array' && schema.items) {
        var itemType = schema.items.$ref
          ? util.capitalize(util.extractModelName(schema.items.$ref))
          : this.translateType(schema.items, imports, context);
        this.addTypeImport(itemType, imports);
        type = this.customSliceHandler(context, schema, itemType);

        if (!type) {
          if (schema.uniqueItems) {
            type = this.setTypeFormat.replace('{type}', itemType);
          } else {
            type = '[]' + itemType;
          }
        }

        type = util.translate(type, this.convertMap);
      } else if (schema.additionalProperties) {
        if (schema.additionalProperties.$ref) {
          var itemType = util.capitalize(util.extractModelName(schema.additionalProperties.$ref));
          this.addTypeImport(itemType, imports);
          type = 'map[string]' + itemType;
        } else {
          type = 'map[string]' + this.translateType(schema.additionalProperties, imports, context);
        }
        type = util.translate(type, this.convertMap);
      } else if (this.convertEnums && schemaType == 'string' && schema.enum) {
        var v = 1;
        var values = _.map(schema.enum, function(value) {
          return {
            name: this.enumName(value),
            value: value,
            valueEscaped: this.escapeGoString(value),
            number: v++
          };
        }.bind(this));

        var enumName = schema['x-enum-name'] || context.propertyName || schema.name;
        if (!enumName) {
          console.log(schema);
        }
        var enumType = this.typeName(enumName) + this.enumSuffix;

        var enumeration = {
          package: this.config.modelPackage || 'io.swagger.model',
          name: enumType,
          varName: this.variableName(enumName),
          values: values,
          hasSlice: false,
          hasSet: false
        };
        var existing = context.enums[enumType];
        if (existing) {
          if (JSON.stringify(existing) != JSON.stringify(enumeration)) {
            console.log('WARNING: Conflicting enumeration "' + enumType + '" detected. ');
          }
        } else {
          context.enums[enumType] = enumeration;
        }

        type = enumType;

        if (schema.required != undefined && !schema.required && !noPointer) {
          type = '*' + type;
        }
      } else {
        type = util.translate(schemaType, this.translateTypeMap, type);
        type = util.translate(schema.format, this.translateFormatMap, type);
        // Translate before and after prepending the pointer asterisk
        type = util.translate(type, this.convertMap);

        var schemaFormat = schema['x-format'] || schema.format;
        type = util.translate(schemaFormat, this.translateFormatMap, type);

        if (schema.required != undefined && !schema.required && !noPointer) {
          type = '*' + type;
        }

        type = util.translate(type, this.convertMap);
        this.addTypeImport(type, imports);
      }
    // schema.$ref needed for response
    } else if (schema.$ref || (schema.schema && schema.schema.$ref)) {
      var $ref = schema.$ref || schema.schema.$ref;
      type = util.capitalize(util.extractModelName($ref));
      var resolved = util.resolveReference(schema.schema || schema, spec, references);

      if (external) {
        type = this.getPackagePrefix(modelPackage) + '.' + type;
        this.addImport(modelPackage, imports);
      }

      if (!noPointer) {
        type = '*' + type;
      }
    } else {
      type = "void";
    }

    return type;
  }

  this.typeName = function(value) {
    return util.capitalize(inflection.camelize(value.replace(/\-/, '_')));
  }

  this.enumName = function(value) {
    value = value.replace(/\-/g, '_');
    value = value.replace(/[^a-zA-Z1-9_ ]/g, '');
    value = value.replace(/\s+/g, ' ');
    value = value.replace(/ /g, '_');
    return util.capitalize(inflection.camelize(value));
  }

  this.variableName = function(value) {
    return inflection.camelize(value.replace(/\-/g, '_'));
  }

  this.methodName = function(value) {
    return inflection.camelize(value.replace(/\-/g, '_'));
  }

  this.escapeGoString = function(value) {
    return '\"' + value
            .replace(/\\/g, '\\')
            .replace(/\r/g, '\\r')
            .replace(/\n/g, '\\n')
            .replace(/\t/g, '\\t')
            .replace(/\"/g, '\\\"')
             + '\"';
  }

  this.updateMaxLengths = function(model, property) {
    model.maxVarNameLength = Math.max(model.maxVarNameLength, property.varname.length);
    model.maxTypeNameLength = Math.max(model.maxTypeNameLength, property.dataType.length);
  }

  this.onPrepare('Context', function(context) {
    context.enums = {};
  });

  this.onDecorate('Operation', function(context) {
    var operation = context.operation;
    var resource = context.resource;

    var okResponse = operation.successResponse || {};
    var returnType = this.translateType(okResponse.schema, resource, context, false) || 'void';
    var returnTypeExternal = this.translateType(okResponse.schema, resource, context, true) || 'void';
    var hasReturn = returnType != 'void';
    var returnDescription = hasReturn ? okResponse.description : null;

    _.assign(operation, {
      methodName : this.methodName(operation['x-resource-operation'] || operation.operationId.replace(/[\s]+/g, '_')),
      globalMethodName : this.methodName(operation.operationId.replace(/[\s]+/g, '_')),
      returnType : returnType,
      returnTypeExternal : returnTypeExternal,
      hasReturn : hasReturn,
      returnDescription : returnDescription
    });
  });

  this.onDecorate('Parameter', function(context) {
    var parameter = context.parameter;
    var resource = context.resource;

    _.assign(parameter, {
      varname : parameter.name,
      dataType : this.translateType(parameter, resource, context, false, false),
      dataTypeExternal : this.translateType(parameter, resource, context, true, false),
      inlineType : this.translateType(parameter, resource, context, false, true),
      inlineTypeExternal : this.translateType(parameter, resource, context, true, true),
      itemType : (parameter.items != null) ? this.translateType(parameter.items, this.translateTypeMap, context) : null
    });
  });

  this.onDecorate('Resource', function(context) {
    var resource = context.resource;

    _.assign(resource, {
      package : this.config.resourcePackage || this.config.package || 'resource',
      structname : util.capitalize(resource.name)
    });
  });

  this.onDecorate('Model', function(context) {
    var model = context.model;
    var name = context.modelName;
    var package = this.config.modelPackage || 'model';

    _.assign(model, {
      package : package,
      structname : util.capitalize(name)
    });

    if (model.references && model.references.length > 0) {
      model.parent = util.capitalize(model.references[0].name);

      _.each(model.references, _.bind(function(reference) {
        if (reference.package && reference.package != package) {
          this.addImport(reference.package, model);
        }
      }, this));
    }
  });

  this.onDecorate('Property', function(context) {
    var property = context.property;
    var model = context.model;

    var imports = [];
    var propertyType = this.translateType(property, imports, context);

    _.each(imports, function(_import) {
      this.addImport(_import, model);
    }.bind(this));

    _.assign(property, {
      tags : {},
      varname : this.variableName(property.name),
      dataType : propertyType,
      imports : imports
    });

    if (property.default != undefined) {
      property.defaultValue = property.dataType == 'string'
        ? this.escapeGoString(property.default)
        : property.default;
    }
  });

  this.onFinalize('Resource', function(context) {
    var resource = context.resource;
    resource.imports = resource.imports || [];
    resource.imports.sort();
  });

  this.onFinalize('Property', function(context) {
    var property = context.property;

    if (!_.isEmpty(property.tags)) {
      var tags = [];
      _.each(property.tags, function(values, tag) {
        tags.push(_.map(values, function(value) {
          if (tag == '_') {
            return value;
          }
          return tag + ':' + this.escapeGoString(value);
        }.bind(this)).join(' '));
      }.bind(this));
      property.tagsString = tags.join(' ');
    }
  });

  this.onFinalize('Model', function(context) {
    var model = context.model;
    model.imports = model.imports || [];
    model.imports.sort();

    model.allImports = _(model.allVars)
      .map(function(item) { return item.imports; })
      .filter(function(item) { return item != null; })
      .flatten()
      .union(model.imports)
      .sort()
      .value();

    var maxVarNameLength = 0;
    var maxTypeNameLength = 0;
    _.each(model.properties, function(property, name) {
      this.updateMaxLengths(model, property);
    }.bind(this));
  });
}
