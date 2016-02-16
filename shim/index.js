var _ref = require('./script/node_modules/nestorbot');
var Robot = _ref.Robot;
var TextMessage = _ref.TextMessage;
var User = _ref.User;
var NestorAdapter = _ref.NestorAdapter;
var Response = _ref.Response;

exports.handle = function(event, ctx) {
  var missingEnv = [];

  for(var requiredEnvProp in event.__nestor_required_env) {
    if(event.__nestor_required_env[requiredEnvProp] == true &&
        (!(requiredEnvProp in event.__nestor_env) ||
         event.__nestor_env[requiredEnvProp] == "")) {
      missingEnv.push(requiredEnvProp);
    }
  }

  var robot = new Robot(relaxEvent.team_uid, relaxEvent.relax_bot_uid, event.__debugMode);

  if(missingEnv.length > 0) {
    var strings = ["You need to set the following environment variables: " + missingEnv.join(', ')];
    var response = new Response(robot);

    response.reply(strings, function() {
      ctx.succeed({
        to_send: robot.toSend
      });
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

    if(relaxEvent.im == true) {
      relaxEvent.text = "<@" + relaxEvent.relax_bot_uid + ">: " + relaxEvent.text;
    }

    message = new TextMessage(user, relaxEvent.text);

    require('script')(robot);

    robot.receive(message, function(done) {
      ctx.succeed({
        to_send: robot.toSend,
      });
    });
  }
}
