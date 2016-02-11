var _ref = require('./script/node_modules/nestorbot');
var Robot = _ref.Robot;
var TextMessage = _ref.TextMessage;
var User = _ref.User;
var NestorAdapter = _ref.NestorAdapter;

exports.handle = function(event, ctx) {
  var missingEnv = [];

  for(var requiredEnvProp in event.__nestor_required_env) {
    if(event.__nestor_required_env[requiredEnvProp] == true &&
        (!(requiredEnvProp in event.__nestor_env) ||
         event.__nestor_env[requiredEnvProp] == "")) {
      missingEnv.push(requiredEnvProp);
    }
  }

  if(missingEnv.length > 0) {
    ctx.succeed({
      to_send: [{strings: ["You need to set the following environment variables: " + missingEnv.join(', ')], reply: true}]
    });
  } else {
    for(var envProp in event.__nestor_env) {
      process.env[envProp] = event.__nestor_env[envProp];
    }

    relaxEvent = event.__relax_event;
    user = new User(relaxEvent.user_uid,
                    {
                      room: relaxEvent.channel_uid
                    });
    message = new TextMessage(user, relaxEvent.text);
    robot = new Robot(relaxEvent.team_uid, relaxEvent.relax_bot_uid, event.__debugMode);

    require('script')(robot);

    robot.receive(message, function(done) {
      ctx.succeed({
        to_send: robot.toSend,
      });
    });
  }
}
