var _ref = require('./script/node_modules/nestorbot');
var Robot = _ref.Robot;
var TextMessage = _ref.TextMessage;
var User = _ref.User;
var NestorAdapter = _ref.NestorAdapter;
var Response = _ref.Response;
var Message = _ref.Message;

exports.handle = function(event, ctx) {
  var missingEnv = [];

  for(var requiredEnvProp in event.__nestor_required_env) {
    if(event.__nestor_required_env[requiredEnvProp] == true &&
        (!(requiredEnvProp in event.__nestor_env) ||
         event.__nestor_env[requiredEnvProp] == "")) {
      missingEnv.push(requiredEnvProp);
    }
  }

  var relaxEvent = event.__relax_event;
  var user = new User(relaxEvent.user_uid,
                    {
                      room: relaxEvent.channel_uid
                    });
  var msg = new TextMessage(user);
  var robot = new Robot(relaxEvent.team_uid, relaxEvent.relax_bot_uid, event.__debugMode);

  for(var envProp in event.__nestor_env) {
    process.env[envProp] = event.__nestor_env[envProp];
  }

  if(missingEnv.length > 0) {
    var strings = ["You need to set the following environment variables: " + missingEnv.join(', ')];
    var response = new Response(robot, msg);

    response.reply(strings, function() {
      ctx.succeed({
        to_send: robot.toSend
      });
    });
  } else {
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
