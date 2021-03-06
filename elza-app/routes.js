//SPDX-License-Identifier: Apache-2.0

var controller = require('./controller.js');

module.exports = function(app){

 app.get('/get_all_groups', function(req, res){
    controller.get_all_groups(req, res);
  });
  app.get('/get_all_tests', function(req, res){
    controller.get_all_tests(req, res);
  });
  app.get('/create_test_group/:generator', function(req, res){
    controller.create_test_group(req, res);
  });
  app.get('/get_test_id/:id', function(req, res){
    controller.get_test_id(req, res);
  });
  app.get('/get_test_student/:name', function(req, res){
      controller.get_test_student(req, res);
  })
  app.get('/prepare_exam/:exam', function(req, res){
      controller.prepare_exam(req, res);
  })
  app.get('/take_test/:exam', function(req, res){
      controller.take_test(req, res);
  })

}
